import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import { Badge } from './ui/badge'
import useStore from '../store/dagStore'
import { Trophy, Zap, Clock, TrendingUp } from 'lucide-react'

function AgentAnalytics() {
  const { agents } = useStore()
  
  // Sort agents by various metrics for leaderboards
  const agentList = Object.values(agents)
  
  const topPerformers = [...agentList]
    .sort((a: any, b: any) => b.tasksCompleted - a.tasksCompleted)
    .slice(0, 5)
  
  const mostEfficient = [...agentList]
    .filter((a: any) => a.tasksCompleted > 0)
    .sort((a: any, b: any) => {
      const aEfficiency = a.avgTokensPerTask || Number.MAX_VALUE
      const bEfficiency = b.avgTokensPerTask || Number.MAX_VALUE
      return aEfficiency - bEfficiency
    })
    .slice(0, 5)

  return (
    <div className="grid gap-4">
      {/* Top Performers */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-sm flex items-center gap-2">
            <Trophy className="h-4 w-4" />
            Top Performers
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {topPerformers.map((agent: any, index) => (
              <div key={agent.name} className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium">{index + 1}.</span>
                  <span className="text-sm">{agent.name}</span>
                  {agent.currentTask && (
                    <Badge variant="secondary" className="text-xs">
                      Active
                    </Badge>
                  )}
                </div>
                <div className="flex items-center gap-4 text-xs text-muted-foreground">
                  <span>{agent.tasksCompleted} tasks</span>
                  {agent.totalTokens > 0 && (
                    <span>{(agent.totalTokens / 1000).toFixed(1)}k tokens</span>
                  )}
                </div>
              </div>
            ))}
            {topPerformers.length === 0 && (
              <p className="text-sm text-muted-foreground">No agents active yet</p>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Most Efficient */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-sm flex items-center gap-2">
            <Zap className="h-4 w-4" />
            Most Efficient
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {mostEfficient.map((agent: any, index) => (
              <div key={agent.name} className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium">{index + 1}.</span>
                  <span className="text-sm">{agent.name}</span>
                </div>
                <div className="text-xs text-muted-foreground">
                  {agent.avgTokensPerTask.toLocaleString()} tokens/task
                </div>
              </div>
            ))}
            {mostEfficient.length === 0 && (
              <p className="text-sm text-muted-foreground">No completed tasks yet</p>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

export default AgentAnalytics