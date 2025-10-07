import { useEffect, useRef, useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import { Badge } from './ui/badge'
import useStore from '../store/dagStore'
import { Clock } from 'lucide-react'

interface TaskBar {
  agent: string
  task: string
  startTime: number
  endTime?: number
  status: string
  tokens: number
}

function AgentGanttChart() {
  const { events, taskStates, agents } = useStore()
  const [taskBars, setTaskBars] = useState<TaskBar[]>([])
  const [timeRange, setTimeRange] = useState({ min: 0, max: 0 })
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    // Build task bars from events and task states
    const bars: TaskBar[] = []
    const agentTasks: Record<string, TaskBar> = {}

    // Process events to build task timeline
    events.forEach((event: any) => {
      if (event.type === 'task') {
        const key = `${event.agent}-${event.task}`
        
        if (event.status === 'started') {
          agentTasks[key] = {
            agent: event.agent,
            task: event.task,
            startTime: event.timestamp,
            status: 'started',
            tokens: 0
          }
        } else if (agentTasks[key]) {
          // Update existing task
          agentTasks[key].status = event.status
          agentTasks[key].tokens = event.tokens || 0
          
          if (event.status === 'completed' || event.status === 'failed') {
            agentTasks[key].endTime = event.timestamp
            bars.push({ ...agentTasks[key] })
            delete agentTasks[key]
          }
        }
      }
    })

    // Add currently running tasks
    Object.values(agentTasks).forEach(task => {
      bars.push({
        ...task,
        endTime: Date.now() / 1000 // Use current time for ongoing tasks
      })
    })

    // Calculate time range
    if (bars.length > 0) {
      const times = bars.flatMap(b => [b.startTime, b.endTime || Date.now() / 1000])
      setTimeRange({
        min: Math.min(...times),
        max: Math.max(...times)
      })
    }

    setTaskBars(bars)
  }, [events, taskStates])

  const getBarPosition = (bar: TaskBar) => {
    if (timeRange.max === timeRange.min) return { left: '0%', width: '100%' }
    
    const duration = timeRange.max - timeRange.min
    const start = ((bar.startTime - timeRange.min) / duration) * 100
    const end = ((bar.endTime || Date.now() / 1000) - timeRange.min) / duration * 100
    
    return {
      left: `${Math.max(0, start)}%`,
      width: `${Math.max(1, end - start)}%`
    }
  }

  const getBarColor = (status: string) => {
    switch (status) {
      case 'completed': return 'bg-green-500'
      case 'failed': return 'bg-red-500'
      case 'started': return 'bg-blue-500'
      default: return 'bg-gray-400'
    }
  }

  const formatTime = (timestamp: number) => {
    const date = new Date(timestamp * 1000)
    return date.toLocaleTimeString()
  }

  // Group bars by agent
  const agentGroups = taskBars.reduce((groups, bar) => {
    if (!groups[bar.agent]) groups[bar.agent] = []
    groups[bar.agent].push(bar)
    return groups
  }, {} as Record<string, TaskBar[]>)

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Clock className="h-5 w-5" />
          Agent Task Timeline
        </CardTitle>
      </CardHeader>
      <CardContent>
        {Object.keys(agentGroups).length > 0 ? (
          <div className="space-y-4">
            {/* Time axis */}
            {timeRange.max > timeRange.min && (
              <div className="flex justify-between text-xs text-muted-foreground px-2">
                <span>{formatTime(timeRange.min)}</span>
                <span>Timeline</span>
                <span>{formatTime(timeRange.max)}</span>
              </div>
            )}
            
            {/* Agent rows */}
            <div className="space-y-3">
              {Object.entries(agentGroups).map(([agent, bars]) => (
                <div key={agent} className="space-y-1">
                  <div className="flex items-center gap-2 mb-1">
                    <span className="text-sm font-medium w-32 truncate">{agent}</span>
                    <Badge variant="outline" className="text-xs">
                      {bars.length} task{bars.length !== 1 ? 's' : ''}
                    </Badge>
                  </div>
                  
                  <div className="relative h-8 bg-muted rounded overflow-hidden" ref={containerRef}>
                    {bars.map((bar, index) => {
                      const position = getBarPosition(bar)
                      return (
                        <div
                          key={`${bar.task}-${index}`}
                          className={`absolute h-6 top-1 ${getBarColor(bar.status)} rounded shadow-sm flex items-center px-1 overflow-hidden`}
                          style={position}
                          title={`${bar.task} (${bar.status})${bar.tokens ? ` - ${bar.tokens} tokens` : ''}`}
                        >
                          <span className="text-xs text-white truncate">
                            {bar.task}
                          </span>
                        </div>
                      )
                    })}
                  </div>
                </div>
              ))}
            </div>
            
            {/* Legend */}
            <div className="flex gap-4 justify-center pt-2 border-t">
              <div className="flex items-center gap-1">
                <div className="w-3 h-3 bg-blue-500 rounded" />
                <span className="text-xs text-muted-foreground">In Progress</span>
              </div>
              <div className="flex items-center gap-1">
                <div className="w-3 h-3 bg-green-500 rounded" />
                <span className="text-xs text-muted-foreground">Completed</span>
              </div>
              <div className="flex items-center gap-1">
                <div className="w-3 h-3 bg-red-500 rounded" />
                <span className="text-xs text-muted-foreground">Failed</span>
              </div>
            </div>
          </div>
        ) : (
          <p className="text-sm text-muted-foreground text-center py-8">
            No task timeline data available yet. Tasks will appear here as agents start working.
          </p>
        )}
      </CardContent>
    </Card>
  )
}

export default AgentGanttChart