import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import { Badge } from './ui/badge'
import { Activity, GitBranch, Clock, DollarSign } from 'lucide-react'

interface AgentStatsProps {
  gitInsights: any
  className?: string
}

export const AgentStats: React.FC<AgentStatsProps> = ({ gitInsights, className }) => {
  if (!gitInsights) {
    return (
      <Card className={className}>
        <CardHeader>
          <CardTitle>Git Insights</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">No git activity yet</p>
        </CardContent>
      </Card>
    )
  }

  const { summary, hotFiles, recentCommits, topContributors, fileTypes } = gitInsights

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Summary Stats */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Activity className="w-5 h-5" />
            Activity Summary
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <p className="text-sm text-muted-foreground">Total Commits</p>
              <p className="text-2xl font-bold">{summary.totalCommits}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Lines Added</p>
              <p className="text-2xl font-bold text-green-500">+{summary.totalLinesAdded}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Lines Removed</p>
              <p className="text-2xl font-bold text-red-500">-{summary.totalLinesRemoved}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Net Lines</p>
              <p className="text-2xl font-bold">{summary.netLines}</p>
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Hot Files */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Hot Files</CardTitle>
          </CardHeader>
          <CardContent>
            {hotFiles && hotFiles.length > 0 ? (
              <div className="space-y-2">
                {hotFiles.slice(0, 5).map((file: any, index: number) => (
                  <div key={file.path} className="flex items-center justify-between text-sm">
                    <span className="font-mono text-xs truncate flex-1">{file.path}</span>
                    <Badge variant="outline">{file.changes}</Badge>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">No file changes yet</p>
            )}
          </CardContent>
        </Card>

        {/* Top Contributors */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Top Contributors</CardTitle>
          </CardHeader>
          <CardContent>
            {topContributors && topContributors.length > 0 ? (
              <div className="space-y-2">
                {topContributors.map((contributor: any, index: number) => (
                  <div key={contributor.agent} className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <div className={`w-2 h-2 rounded-full ${
                        index === 0 ? 'bg-yellow-500' :
                        index === 1 ? 'bg-gray-400' :
                        index === 2 ? 'bg-orange-600' : 'bg-muted'
                      }`} />
                      <span className="text-sm font-medium">{contributor.agent}</span>
                    </div>
                    <Badge>{contributor.commits} commits</Badge>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">No commits yet</p>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Recent Commits */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <GitBranch className="w-5 h-5" />
            Recent Commits
          </CardTitle>
        </CardHeader>
        <CardContent>
          {recentCommits && recentCommits.length > 0 ? (
            <div className="space-y-3">
              {recentCommits.slice(0, 5).map((commit: any) => (
                <div key={commit.sha} className="flex items-start gap-3 p-2 rounded-lg border bg-muted/20">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="font-mono text-xs text-muted-foreground">
                        {commit.sha.substring(0, 7)}
                      </span>
                      <Badge variant="outline" className="text-xs">
                        {commit.agent}
                      </Badge>
                      <span className="text-xs text-muted-foreground">
                        {new Date(commit.timestamp * 1000).toLocaleString()}
                      </span>
                    </div>
                    <p className="text-sm">{commit.message || 'No message'}</p>
                    <div className="flex gap-4 mt-1 text-xs text-muted-foreground">
                      <span className="text-green-600">+{commit.linesAdded}</span>
                      <span className="text-red-600">-{commit.linesRemoved}</span>
                      <span>{commit.filesChanged} files</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-sm text-muted-foreground">No commits yet</p>
          )}
        </CardContent>
      </Card>

      {/* File Types */}
      {fileTypes && Object.keys(fileTypes).length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-base">File Types Changed</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex gap-2 flex-wrap">
              {Object.entries(fileTypes).map(([ext, count]) => (
                <Badge key={ext} variant="secondary">
                  .{ext}: {count}
                </Badge>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}