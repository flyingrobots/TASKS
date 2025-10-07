import { useEffect, useState, useMemo } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from './components/ui/card'
import { Badge } from './components/ui/badge'
import { Progress } from './components/ui/progress'
import { Sheet, SheetContent, SheetHeader, SheetTitle } from './components/ui/sheet'
import { Toaster } from './components/ui/sonner'
import { toast } from 'sonner'
import { DAGViewer } from './components/DAGViewer'
import AgentAnalytics from './components/AgentAnalytics'
import GitAnalytics from './components/GitAnalytics'
import AgentGanttChart from './components/AgentGanttChart'
import { AgentEventStream } from './components/AgentEventStream'
import { GitEventStream } from './components/GitEventStream'
import { useWebSocket } from './hooks/useWebSocket'
import useStore from './store/dagStore'
import { loadConfig, getWebSocketUrl } from './config'
import { 
  Activity, 
  GitBranch, 
  Users, 
  Clock,
  // TrendingUp,
  CheckCircle2,
  XCircle,
  AlertCircle
} from 'lucide-react'

function App() {
  const { dagData, taskStates, agents, gitStats, handleWebSocketMessage, setDagData, setWsConnected } = useStore()
  const [taskMetadata, setTaskMetadata] = useState<Record<string, any>>({})
  const [wsUrl, setWsUrl] = useState<string>('ws://localhost:8080')
  
  // Load config on mount
  useEffect(() => {
    loadConfig().then(() => {
      setWsUrl(getWebSocketUrl())
    })
  }, [])
  
  const { isConnected } = useWebSocket({
    url: wsUrl,
    onMessage: handleWebSocketMessage,
    onOpen: () => {
      setWsConnected(true)
      toast.success('Connected to live updates', {
        description: 'Real-time task and git activity enabled',
        duration: 3000,
      })
    },
    onClose: () => {
      setWsConnected(false)
      toast.error('Disconnected from live updates', {
        description: 'Attempting to reconnect...',
        duration: 3000,
      })
    }
  })
  const [startTime, setStartTime] = useState<number | null>(null)
  const [elapsedTime, setElapsedTime] = useState(0)
  const [showDashboard, setShowDashboard] = useState(false)
  
  // Load DAG data from backend state endpoint
  useEffect(() => {
    if (!dagData) {
      fetch('http://localhost:8080/state')
        .then(res => res.json())
        .then(state => {
          if (state.tasksMetadata && state.tasksMetadata.tasks) {
            // Store task metadata for lookup
            const taskLookup: Record<string, any> = {}
            state.tasksMetadata.tasks.forEach((task: any) => {
              taskLookup[task.id] = task
            })
            setTaskMetadata(taskLookup)
            
            // Build DAG from tasks.json format
            const topo_order = state.tasksMetadata.tasks.map((t: any) => t.id)
            const dependencies = state.tasksMetadata.dependencies || []
            
            const dagFromTasks = {
              ...state.tasksMetadata,
              topo_order,
              reduced_edges_sample: dependencies.map((d: any) => [d.from, d.to]),
              features: state.features?.features || [],
              waves: state.waves?.waves || [],
              metrics: state.tasksMetadata.generated || state.dag?.metrics || {}
            }
            setDagData(dagFromTasks)
          } else if (state.dag && state.dag.topo_order) {
            // Fall back to original dag.json format
            setDagData(state.dag)
          }
        })
        .catch(err => console.error('Failed to load DAG data from backend:', err))
    }
  }, [dagData, setDagData])

  // Track start time and elapsed time
  useEffect(() => {
    const hasStartedTasks = Object.values(taskStates).some((t: any) => 
      t.status !== 'pending' && t.status !== 'blocked'
    )
    
    if (hasStartedTasks && !startTime) {
      // Use current time as start time if tasks don't have timestamps
      setStartTime(Date.now() / 1000)
    }
  }, [taskStates, startTime])

  useEffect(() => {
    if (!startTime) return
    
    // Check if all tasks are completed or failed
    const allTasksDone = Object.values(taskStates).length > 0 && 
      Object.values(taskStates).every((t: any) => 
        t.status === 'completed' || t.status === 'failed'
      )
    
    // Don't run timer if all tasks are done
    if (allTasksDone) return
    
    const interval = setInterval(() => {
      setElapsedTime(Date.now() / 1000 - startTime)
    }, 1000)
    
    return () => clearInterval(interval)
  }, [startTime, taskStates])

  const formatTime = (seconds: number) => {
    const hours = Math.floor(seconds / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    const secs = Math.floor(seconds % 60)
    
    if (hours > 0) {
      return `${hours}h ${minutes}m ${secs}s`
    } else if (minutes > 0) {
      return `${minutes}m ${secs}s`
    }
    return `${secs}s`
  }

  const progress = useMemo(() => {
    const total = dagData?.topo_order?.length || Object.keys(taskStates).length || 0
    const inProgressOrCompleted = Object.values(taskStates).filter((t: any) => 
      t.status === 'in_progress' || t.status === 'completed' || t.status === 'started'
    ).length
    const percentage = total > 0 ? (inProgressOrCompleted / total) * 100 : 0
    return { completed: inProgressOrCompleted, total, percentage }
  }, [taskStates, dagData])

  const taskStats = useMemo(() => {
    const stats = {
      pending: 0,
      started: 0,
      completed: 0,
      failed: 0,
      blocked: 0
    }
    
    Object.values(taskStates).forEach((task: any) => {
      const status = task.status || 'pending'
      if (status in stats) {
        stats[status as keyof typeof stats]++
      }
    })
    
    return stats
  }, [taskStates])

  const agentStats = useMemo(() => {
    const activeAgents = Object.values(agents).filter((agent: any) => agent.currentTask).length
    const totalTasks = Object.values(agents).reduce((sum, agent: any) => 
      sum + agent.tasksCompleted + agent.tasksFailed + agent.tasksStarted, 0
    )
    const totalTokens = Object.values(agents).reduce((sum, agent: any) => 
      sum + agent.totalTokens, 0
    )
    
    return { activeAgents, totalTasks, totalTokens }
  }, [agents])

  return (
    <div className="fixed inset-0 overflow-hidden bg-background">
      {/* Fullscreen DAG Graph */}
      <DAGViewer 
        dagData={dagData} 
        taskStates={taskStates} 
        taskMetadata={taskMetadata}
        onDashboardToggle={setShowDashboard}
      />

      {/* Toast notifications */}
      <Toaster position="bottom-right" />

      {/* Transparent Header Overlay */}
      <header className="absolute top-0 left-0 right-0 z-20 bg-background/80 backdrop-blur-sm border-b">
        <div className="container mx-auto px-4 py-3">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-xl font-bold">DAG Task Analytics</h1>
              <p className="text-xs text-muted-foreground">Real-time task execution and agent performance</p>
            </div>
            <div className="flex items-center gap-3">
              <Badge variant={isConnected ? 'default' : 'destructive'} className="text-xs">
                {isConnected ? 'Live' : 'Disconnected'}
              </Badge>
              {progress.total > 0 && (
                <div className="flex items-center gap-2">
                  <Progress value={progress.percentage} className="w-24 h-2" />
                  <span className="text-xs font-medium">
                    {progress.completed}/{progress.total}
                  </span>
                </div>
              )}
              {elapsedTime > 0 && (
                <div className="flex items-center gap-1">
                  <Clock className="h-3 w-3 text-muted-foreground" />
                  <span className="text-xs font-medium">{formatTime(elapsedTime)}</span>
                </div>
              )}
            </div>
          </div>
        </div>
      </header>

      {/* Dashboard Drawer */}
      <Sheet open={showDashboard} onOpenChange={setShowDashboard}>
        <SheetContent side="bottom" className="h-[80vh] overflow-y-auto bg-white dark:bg-gray-950 border-t">
          <SheetHeader>
            <SheetTitle>Analytics Dashboard</SheetTitle>
          </SheetHeader>
          <div className="container mx-auto px-4 py-6">
            <div className="grid gap-6">
              {/* Agent Gantt Chart */}
              <AgentGanttChart />

              {/* Agent and Git Analytics Side by Side */}
              <div className="grid gap-6 lg:grid-cols-2">
                {/* Agent Analytics */}
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <h2 className="text-xl font-semibold flex items-center gap-2">
                      <Users className="h-5 w-5" />
                      Agent Performance
                    </h2>
                    <div className="flex gap-2">
                      <Badge variant="outline">
                        {agentStats.activeAgents} Active
                      </Badge>
                      <Badge variant="outline">
                        {Object.keys(agents).length} Total
                      </Badge>
                    </div>
                  </div>
              
              <div className="grid gap-4 md:grid-cols-3">
                <Card>
                  <CardHeader className="pb-2">
                    <CardTitle className="text-sm">Total Tasks</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{agentStats.totalTasks}</div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardHeader className="pb-2">
                    <CardTitle className="text-sm">Total Tokens</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">
                      {agentStats.totalTokens.toLocaleString()}
                    </div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardHeader className="pb-2">
                    <CardTitle className="text-sm">Avg Tokens/Task</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">
                      {agentStats.totalTasks > 0 
                        ? Math.round(agentStats.totalTokens / agentStats.totalTasks).toLocaleString()
                        : 0}
                    </div>
                  </CardContent>
                </Card>
              </div>
              
              <AgentAnalytics />
            </div>

            {/* Git Analytics */}
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <h2 className="text-xl font-semibold flex items-center gap-2">
                  <GitBranch className="h-5 w-5" />
                  Git Activity
                </h2>
                <Badge variant="outline">
                  {gitStats.totalCommits} Commits
                </Badge>
              </div>
              
              <div className="grid gap-4 md:grid-cols-3">
                <Card>
                  <CardHeader className="pb-2">
                    <CardTitle className="text-sm">Lines Added</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold text-green-600">
                      +{gitStats.totalLinesAdded.toLocaleString()}
                    </div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardHeader className="pb-2">
                    <CardTitle className="text-sm">Lines Removed</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold text-red-600">
                      -{gitStats.totalLinesRemoved.toLocaleString()}
                    </div>
                  </CardContent>
                </Card>
                
                <Card>
                  <CardHeader className="pb-2">
                    <CardTitle className="text-sm">Files Changed</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">
                      {Object.keys(gitStats.fileChangeFrequency || {}).length}
                    </div>
                  </CardContent>
                </Card>
              </div>
              
              <GitAnalytics />
            </div>
          </div>

              {/* Event Streams Side by Side */}
              <div className="grid gap-6 lg:grid-cols-2">
                <AgentEventStream />
                <GitEventStream />
              </div>
            </div>
          </div>
        </SheetContent>
      </Sheet>
    </div>
  )
}

export default App