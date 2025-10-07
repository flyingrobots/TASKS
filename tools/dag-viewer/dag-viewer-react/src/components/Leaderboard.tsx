import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import { Badge } from './ui/badge'
import { Trophy, TrendingUp, GitCommit, Zap } from 'lucide-react'

interface LeaderboardProps {
  agents: any[]
  className?: string
}

export const Leaderboard: React.FC<LeaderboardProps> = ({ agents, className }) => {
  const getScoreBadgeVariant = (score: number) => {
    if (score >= 80) return 'default'
    if (score >= 60) return 'secondary'
    if (score >= 40) return 'outline'
    return 'destructive'
  }

  const getRankIcon = (rank: number) => {
    switch (rank) {
      case 1:
        return <Trophy className="w-5 h-5 text-yellow-500" />
      case 2:
        return <Trophy className="w-5 h-5 text-gray-400" />
      case 3:
        return <Trophy className="w-5 h-5 text-orange-600" />
      default:
        return <span className="w-5 h-5 flex items-center justify-center text-muted-foreground">#{rank}</span>
    }
  }

  return (
    <Card className={className}>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Trophy className="w-5 h-5" />
          Agent Leaderboard
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {agents.length === 0 ? (
            <p className="text-muted-foreground text-center py-8">No agents active yet</p>
          ) : (
            agents.map((agent, index) => (
              <div key={agent.name} className="flex items-start gap-3 p-3 rounded-lg border bg-card">
                <div className="flex-shrink-0 mt-1">
                  {getRankIcon(index + 1)}
                </div>
                
                <div className="flex-1 space-y-2">
                  <div className="flex items-center justify-between">
                    <h4 className="font-semibold">{agent.name}</h4>
                    <Badge variant={getScoreBadgeVariant(agent.productivityScore || 0)}>
                      Score: {agent.productivityScore || 0}
                    </Badge>
                  </div>
                  
                  <div className="grid grid-cols-2 gap-2 text-sm">
                    <div className="flex items-center gap-1">
                      <TrendingUp className="w-3 h-3 text-green-500" />
                      <span className="text-muted-foreground">Tasks:</span>
                      <span className="font-medium">{agent.tasksCompleted || 0}</span>
                      {agent.tasksFailed > 0 && (
                        <span className="text-red-500">({agent.tasksFailed} failed)</span>
                      )}
                    </div>
                    
                    <div className="flex items-center gap-1">
                      <GitCommit className="w-3 h-3 text-blue-500" />
                      <span className="text-muted-foreground">Commits:</span>
                      <span className="font-medium">{agent.commits || 0}</span>
                    </div>
                    
                    <div className="flex items-center gap-1">
                      <Zap className="w-3 h-3 text-yellow-500" />
                      <span className="text-muted-foreground">Tokens:</span>
                      <span className="font-medium">{agent.totalTokensK || '0k'}</span>
                    </div>
                    
                    <div className="flex items-center gap-1">
                      <span className="text-muted-foreground">Cost:</span>
                      <span className="font-medium">{agent.tokenCost || '$0'}</span>
                    </div>
                  </div>
                  
                  {agent.currentTask && (
                    <div className="flex items-center gap-1 text-sm">
                      <div className="w-2 h-2 bg-blue-500 rounded-full animate-pulse" />
                      <span className="text-muted-foreground">Working on:</span>
                      <span className="font-mono text-xs">{agent.currentTask}</span>
                    </div>
                  )}
                  
                  <div className="flex gap-2 flex-wrap">
                    {agent.successRate && (
                      <Badge variant="outline" className="text-xs">
                        {agent.successRate}% success
                      </Badge>
                    )}
                    {agent.avgCommitSize && parseInt(agent.avgCommitSize) > 0 && (
                      <Badge variant="outline" className="text-xs">
                        ~{agent.avgCommitSize} lines/commit
                      </Badge>
                    )}
                    {agent.costPerTask && (
                      <Badge variant="outline" className="text-xs">
                        {agent.costPerTask}/task
                      </Badge>
                    )}
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </CardContent>
    </Card>
  )
}