import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import { Badge } from './ui/badge'
import useStore from '../store/dagStore'
import { FileText, GitCommit, Code } from 'lucide-react'

function GitAnalytics() {
  const { gitStats, agents } = useStore()
  
  // Get top committers
  const topCommitters = Object.entries(agents)
    .filter(([_, agent]: [string, any]) => agent.commits > 0)
    .sort(([_, a]: [string, any], [__, b]: [string, any]) => b.commits - a.commits)
    .slice(0, 5)
  
  // Get hot files
  const hotFiles = Object.entries(gitStats.fileChangeFrequency || {})
    .sort(([_, a]: [string, any], [__, b]: [string, any]) => b - a)
    .slice(0, 5)

  return (
    <div className="grid gap-4">
      {/* Top Committers */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-sm flex items-center gap-2">
            <GitCommit className="h-4 w-4" />
            Top Committers
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {topCommitters.map(([name, agent]: [string, any], index) => (
              <div key={name} className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium">{index + 1}.</span>
                  <span className="text-sm">{name}</span>
                </div>
                <div className="flex items-center gap-4 text-xs text-muted-foreground">
                  <span>{agent.commits} commits</span>
                  <span className="text-green-600">+{agent.linesAdded}</span>
                  <span className="text-red-600">-{agent.linesRemoved}</span>
                </div>
              </div>
            ))}
            {topCommitters.length === 0 && (
              <p className="text-sm text-muted-foreground">No commits yet</p>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Hot Files */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-sm flex items-center gap-2">
            <FileText className="h-4 w-4" />
            Hot Files
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {hotFiles.map(([file, changes]: [string, any], index) => (
              <div key={file} className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-medium">{index + 1}.</span>
                  <span className="text-sm truncate max-w-[200px]" title={file}>
                    {file.split('/').pop()}
                  </span>
                </div>
                <Badge variant="secondary" className="text-xs">
                  {changes} changes
                </Badge>
              </div>
            ))}
            {hotFiles.length === 0 && (
              <p className="text-sm text-muted-foreground">No file changes yet</p>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

export default GitAnalytics