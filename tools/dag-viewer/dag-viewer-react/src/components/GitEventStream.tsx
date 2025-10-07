import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import { Badge } from './ui/badge'
import { ScrollArea } from './ui/scroll-area'
import { 
  GitCommit, 
  GitBranch, 
  GitMerge,
  GitPullRequest,
  FileText,
  Plus,
  Minus,
  Edit,
  User,
  Clock,
  Hash
} from 'lucide-react'
import useStore from '../store/dagStore'
import { motion, AnimatePresence } from 'framer-motion'

interface GitEvent {
  id: string
  type: 'commit' | 'file-change' | 'branch' | 'merge'
  author: string
  timestamp: number
  hash?: string
  message?: string
  filesChanged?: number
  linesAdded?: number
  linesRemoved?: number
  files?: string[]
  branch?: string
}

const eventIcons = {
  'commit': <GitCommit className="h-4 w-4 text-blue-500" />,
  'file-change': <FileText className="h-4 w-4 text-purple-500" />,
  'branch': <GitBranch className="h-4 w-4 text-green-500" />,
  'merge': <GitMerge className="h-4 w-4 text-orange-500" />
}

const eventColors = {
  'commit': 'bg-blue-500/10 text-blue-700 border-blue-500/30',
  'file-change': 'bg-purple-500/10 text-purple-700 border-purple-500/30',
  'branch': 'bg-green-500/10 text-green-700 border-green-500/30',
  'merge': 'bg-orange-500/10 text-orange-700 border-orange-500/30'
}

export const GitEventStream: React.FC = () => {
  const { events, agents, gitInsights } = useStore()
  const [gitEvents, setGitEvents] = useState<GitEvent[]>([])
  const [filter, setFilter] = useState<string>('all')

  useEffect(() => {
    // Convert git-related events to git events
    const processedEvents: GitEvent[] = []
    
    // Process events from the event stream
    events.forEach((event: any, index: number) => {
      if (event.type === 'git-commit' || event.type === 'commit') {
        processedEvents.push({
          id: `git-${index}-${event.timestamp || Date.now()}`,
          type: 'commit',
          author: event.author || event.agent || 'Unknown',
          timestamp: event.timestamp || Date.now(),
          hash: event.hash || event.sha,
          message: event.message,
          filesChanged: event.filesChanged,
          linesAdded: event.linesAdded,
          linesRemoved: event.linesRemoved,
          files: event.files
        })
      }
    })

    // Process git insights if available
    if (gitInsights?.recentCommits) {
      gitInsights.recentCommits.forEach((commit: any) => {
        processedEvents.push({
          id: `insight-${commit.sha}`,
          type: 'commit',
          author: commit.author || 'Unknown',
          timestamp: commit.timestamp || Date.now(),
          hash: commit.sha,
          message: commit.message,
          filesChanged: commit.filesChanged,
          linesAdded: commit.linesAdded,
          linesRemoved: commit.linesRemoved
        })
      })
    }

    // Add synthetic events from agent commits
    Object.entries(agents).forEach(([agentName, agentData]) => {
      if (agentData.lastCommit) {
        processedEvents.push({
          id: `agent-commit-${agentName}-${agentData.lastCommit.sha}`,
          type: 'commit',
          author: agentName,
          timestamp: agentData.lastCommit.timestamp || agentData.lastSeen || Date.now(),
          hash: agentData.lastCommit.sha,
          message: agentData.lastCommit.message,
          filesChanged: agentData.filesChangedCount,
          linesAdded: agentData.linesAdded,
          linesRemoved: agentData.linesRemoved
        })
      }
    })

    // Remove duplicates based on hash
    const uniqueEvents = processedEvents.filter((event, index, self) => 
      !event.hash || index === self.findIndex(e => e.hash === event.hash)
    )

    // Sort by timestamp (newest first)
    uniqueEvents.sort((a, b) => b.timestamp - a.timestamp)
    
    // Limit to last 50 events
    setGitEvents(uniqueEvents.slice(0, 50))
  }, [events, agents, gitInsights])

  const filteredEvents = filter === 'all' 
    ? gitEvents 
    : gitEvents.filter(e => e.author === filter)

  const uniqueAuthors = Array.from(new Set(gitEvents.map(e => e.author)))

  const formatTime = (timestamp: number) => {
    const date = new Date(timestamp)
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    
    if (diff < 60000) return 'just now'
    if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`
    
    return date.toLocaleTimeString()
  }

  const formatHash = (hash?: string) => {
    return hash ? hash.substring(0, 7) : '-------'
  }

  return (
    <Card className="overflow-hidden">
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2">
            <GitBranch className="h-5 w-5" />
            Git Event Stream
          </CardTitle>
          <div className="flex gap-2">
            <Badge 
              variant={filter === 'all' ? 'default' : 'outline'}
              className="cursor-pointer"
              onClick={() => setFilter('all')}
            >
              All
            </Badge>
            {uniqueAuthors.slice(0, 3).map(author => (
              <Badge 
                key={author}
                variant={filter === author ? 'default' : 'outline'}
                className="cursor-pointer"
                onClick={() => setFilter(author)}
              >
                {author}
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
                        <span className="font-medium text-sm">{event.author}</span>
                        {event.hash && (
                          <Badge variant="outline" className="text-xs font-mono">
                            {formatHash(event.hash)}
                          </Badge>
                        )}
                      </div>
                      
                      {event.message && (
                        <p className="text-sm text-muted-foreground mt-1 line-clamp-2">
                          {event.message}
                        </p>
                      )}
                      
                      <div className="flex items-center gap-3 mt-2 flex-wrap">
                        <div className="flex items-center gap-1">
                          <Clock className="h-3 w-3 text-muted-foreground" />
                          <span className="text-xs text-muted-foreground">
                            {formatTime(event.timestamp)}
                          </span>
                        </div>
                        
                        {event.filesChanged !== undefined && (
                          <div className="flex items-center gap-1">
                            <FileText className="h-3 w-3 text-muted-foreground" />
                            <span className="text-xs text-muted-foreground">
                              {event.filesChanged} file{event.filesChanged !== 1 ? 's' : ''}
                            </span>
                          </div>
                        )}
                        
                        {event.linesAdded !== undefined && event.linesAdded > 0 && (
                          <div className="flex items-center gap-1">
                            <Plus className="h-3 w-3 text-green-500" />
                            <span className="text-xs text-green-600">
                              +{event.linesAdded}
                            </span>
                          </div>
                        )}
                        
                        {event.linesRemoved !== undefined && event.linesRemoved > 0 && (
                          <div className="flex items-center gap-1">
                            <Minus className="h-3 w-3 text-red-500" />
                            <span className="text-xs text-red-600">
                              -{event.linesRemoved}
                            </span>
                          </div>
                        )}
                      </div>
                      
                      {event.files && event.files.length > 0 && (
                        <div className="mt-2 pt-2 border-t">
                          <div className="flex flex-wrap gap-1">
                            {event.files.slice(0, 3).map((file, i) => (
                              <Badge key={i} variant="secondary" className="text-xs">
                                {file.split('/').pop()}
                              </Badge>
                            ))}
                            {event.files.length > 3 && (
                              <Badge variant="secondary" className="text-xs">
                                +{event.files.length - 3} more
                              </Badge>
                            )}
                          </div>
                        </div>
                      )}
                    </div>
                    
                    {/* Live indicator for recent events */}
                    {index === 0 && Date.now() - event.timestamp < 10000 && (
                      <div className="absolute top-2 right-2">
                        <div className="flex items-center gap-1">
                          <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
                          <span className="text-xs text-green-600">New</span>
                        </div>
                      </div>
                    )}
                  </div>
                </motion.div>
              ))}
            </AnimatePresence>
            
            {filteredEvents.length === 0 && (
              <div className="text-center py-8 text-muted-foreground">
                <GitCommit className="h-8 w-8 mx-auto mb-2 opacity-50" />
                <p className="text-sm">No git events yet</p>
              </div>
            )}
          </div>
        </ScrollArea>
      </CardContent>
    </Card>
  )
}