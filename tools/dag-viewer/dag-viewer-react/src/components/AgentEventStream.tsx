import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import { Badge } from './ui/badge'
import { ScrollArea } from './ui/scroll-area'
import { 
  Activity, 
  CheckCircle2, 
  XCircle, 
  AlertCircle,
  Clock,
  Zap,
  GitCommit,
  User,
  Hash
} from 'lucide-react'
import useStore from '../store/dagStore'
import { motion, AnimatePresence } from 'framer-motion'

interface AgentEvent {
  id: string
  type: 'task-started' | 'task-completed' | 'task-failed' | 'agent-active' | 'tokens-spent'
  agent: string
  task?: string
  timestamp: number
  status?: string
  tokens?: number
  duration?: number
  message?: string
}

const eventIcons = {
  'task-started': <Activity className="h-4 w-4 text-blue-500" />,
  'task-completed': <CheckCircle2 className="h-4 w-4 text-green-500" />,
  'task-failed': <XCircle className="h-4 w-4 text-red-500" />,
  'agent-active': <User className="h-4 w-4 text-purple-500" />,
  'tokens-spent': <Zap className="h-4 w-4 text-yellow-500" />
}

const eventColors = {
  'task-started': 'bg-blue-500/10 text-blue-700 border-blue-500/30',
  'task-completed': 'bg-green-500/10 text-green-700 border-green-500/30',
  'task-failed': 'bg-red-500/10 text-red-700 border-red-500/30',
  'agent-active': 'bg-purple-500/10 text-purple-700 border-purple-500/30',
  'tokens-spent': 'bg-yellow-500/10 text-yellow-700 border-yellow-500/30'
}

export const AgentEventStream: React.FC = () => {
  const { events, taskStates, agents } = useStore()
  const [agentEvents, setAgentEvents] = useState<AgentEvent[]>([])
  const [filter, setFilter] = useState<string>('all')

  useEffect(() => {
    // Convert events and task state changes to agent events
    const processedEvents: AgentEvent[] = []
    
    // Process regular events
    events.forEach((event: any, index: number) => {
      if (event.type === 'task' && event.agent) {
        processedEvents.push({
          id: `event-${index}-${event.timestamp || Date.now()}`,
          type: event.status === 'completed' ? 'task-completed' : 
                event.status === 'failed' ? 'task-failed' : 'task-started',
          agent: event.agent,
          task: event.task,
          timestamp: event.timestamp || Date.now(),
          status: event.status,
          tokens: event.tokens,
          duration: event.duration
        })
      }
    })

    // Process task state changes
    Object.entries(taskStates).forEach(([taskId, state]: [string, any]) => {
      if (state.agent && state.timestamp) {
        const eventType = state.status === 'completed' ? 'task-completed' :
                         state.status === 'failed' ? 'task-failed' : 
                         state.status === 'started' ? 'task-started' : null
        
        if (eventType) {
          processedEvents.push({
            id: `task-${taskId}-${state.timestamp}`,
            type: eventType,
            agent: state.agent,
            task: taskId,
            timestamp: state.timestamp * 1000, // Convert to milliseconds
            status: state.status,
            tokens: state.tokensSpent
          })
        }
      }
    })

    // Sort by timestamp (newest first)
    processedEvents.sort((a, b) => b.timestamp - a.timestamp)
    
    // Limit to last 50 events
    setAgentEvents(processedEvents.slice(0, 50))
  }, [events, taskStates])

  const filteredEvents = filter === 'all' 
    ? agentEvents 
    : agentEvents.filter(e => e.agent === filter)

  const uniqueAgents = Array.from(new Set(agentEvents.map(e => e.agent)))

  const formatTime = (timestamp: number) => {
    const date = new Date(timestamp)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    
    if (diff < 60000) return 'just now'
    if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`
    
    return date.toLocaleTimeString()
  }

  return (
    <Card className="overflow-hidden">
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2">
            <Activity className="h-5 w-5" />
            Agent Event Stream
          </CardTitle>
          <div className="flex gap-2">
            <Badge 
              variant={filter === 'all' ? 'default' : 'outline'}
              className="cursor-pointer"
              onClick={() => setFilter('all')}
            >
              All
            </Badge>
            {uniqueAgents.slice(0, 3).map(agent => (
              <Badge 
                key={agent}
                variant={filter === agent ? 'default' : 'outline'}
                className="cursor-pointer"
                onClick={() => setFilter(agent)}
              >
                {agent}
              </Badge>
            ))}
          </div>
        </div>
      </CardHeader>
      <CardContent className="p-0">
        <ScrollArea className="h-[500px]">
          <div className="p-4 space-y-2">
            <AnimatePresence initial={false}>
              {filteredEvents.map((event, index) => (
                <motion.div
                  key={event.id}
                  initial={{ opacity: 0, y: -20, scale: 0.95 }}
                  animate={{ opacity: 1, y: 0, scale: 1 }}
                  exit={{ opacity: 0, x: -100 }}
                  transition={{ 
                    duration: 0.3,
                    delay: index * 0.05,
                    type: "spring",
                    stiffness: 200
                  }}
                  className="group"
                >
                  <div className={`
                    relative flex items-start gap-3 p-3 rounded-lg border 
                    ${eventColors[event.type]}
                    transition-all duration-200 hover:scale-[1.02] hover:shadow-md
                  `}>
                    {/* Icon */}
                    <div className="mt-0.5">
                      {eventIcons[event.type]}
                    </div>
                    
                    {/* Content */}
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 flex-wrap">
                        <span className="font-medium text-sm">{event.agent}</span>
                        {event.type === 'task-completed' && (
                          <Badge variant="default" className="text-xs">
                            Completed
                          </Badge>
                        )}
                        {event.type === 'task-failed' && (
                          <Badge variant="destructive" className="text-xs">
                            Failed
                          </Badge>
                        )}
                        {event.type === 'task-started' && (
                          <Badge variant="secondary" className="text-xs">
                            Started
                          </Badge>
                        )}
                      </div>
                      
                      {event.task && (
                        <div className="flex items-center gap-2 mt-1">
                          <Hash className="h-3 w-3 text-muted-foreground" />
                          <span className="text-sm text-muted-foreground truncate">
                            {event.task}
                          </span>
                        </div>
                      )}
                      
                      <div className="flex items-center gap-3 mt-2">
                        <div className="flex items-center gap-1">
                          <Clock className="h-3 w-3 text-muted-foreground" />
                          <span className="text-xs text-muted-foreground">
                            {formatTime(event.timestamp)}
                          </span>
                        </div>
                        
                        {event.tokens && (
                          <div className="flex items-center gap-1">
                            <Zap className="h-3 w-3 text-yellow-500" />
                            <span className="text-xs text-muted-foreground">
                              {event.tokens.toLocaleString()} tokens
                            </span>
                          </div>
                        )}
                        
                        {event.duration && (
                          <div className="flex items-center gap-1">
                            <Clock className="h-3 w-3 text-muted-foreground" />
                            <span className="text-xs text-muted-foreground">
                              {(event.duration / 1000).toFixed(1)}s
                            </span>
                          </div>
                        )}
                      </div>
                    </div>
                    
                    {/* Live indicator for recent events */}
                    {index === 0 && Date.now() - event.timestamp < 10000 && (
                      <div className="absolute top-2 right-2">
                        <div className="flex items-center gap-1">
                          <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
                          <span className="text-xs text-green-600">Live</span>
                        </div>
                      </div>
                    )}
                  </div>
                </motion.div>
              ))}
            </AnimatePresence>
            
            {filteredEvents.length === 0 && (
              <div className="text-center py-8 text-muted-foreground">
                <AlertCircle className="h-8 w-8 mx-auto mb-2 opacity-50" />
                <p className="text-sm">No agent events yet</p>
              </div>
            )}
          </div>
        </ScrollArea>
      </CardContent>
    </Card>
  )
}