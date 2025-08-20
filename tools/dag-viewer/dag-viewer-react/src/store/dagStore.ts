import { create } from 'zustand'
import { toast } from 'sonner'

interface TaskState {
  status: string
  agent?: string
  timestamp?: number
  startTime?: number
  tokensSpent?: number
  startTokens?: number
}

interface AgentData {
  name: string
  tasksCompleted: number
  tasksFailed: number
  tasksStarted: number
  taskTimes: number[]
  currentTask: string | null
  totalTokens: number
  tokensPerTask: number[]
  avgTokensPerTask: number
  tokenEfficiency: number
  commits: number
  linesAdded: number
  linesRemoved: number
  filesChanged: Set<string>
  commitSizes: number[]
  lastCommit: any | null
  tasksPerCommit: number
  linesPerTask: number
  commitFrequency: number
  tokensPerLine: number
  firstSeen: number
  lastSeen: number
  activeTime: number
  path: any[]
  productivityScore?: number
  successRate?: string
  avgTaskTime?: number
  medianTaskTime?: number
  totalTokensK?: string
  tokenCost?: string
  costPerTask?: string
  avgCommitSize?: string
  codeChurn?: number
  netLines?: number
  filesChangedCount?: number
}

interface GitInsights {
  summary: {
    totalCommits: number
    totalLinesAdded: number
    totalLinesRemoved: number
    netLines: number
    avgCommitSize: string
  }
  hotFiles: Array<{ path: string; changes: number }>
  recentCommits: any[]
  hourlyActivity: number[]
  topContributors: Array<{ agent: string; commits: number }>
  fileTypes: Record<string, number>
}

interface GitStats {
  totalCommits: number
  totalLinesAdded: number
  totalLinesRemoved: number
  fileChangeFrequency: Record<string, number>
}

interface DagState {
  dagData: any | null
  taskStates: Record<string, TaskState>
  agents: Record<string, AgentData>
  leaderboard: AgentData[]
  gitInsights: GitInsights | null
  gitStats: GitStats
  events: any[]
  wsConnected: boolean
  
  // Actions
  setDagData: (data: any) => void
  updateTaskState: (taskId: string, state: TaskState) => void
  updateTaskStates: (states: Record<string, TaskState>) => void
  updateAgent: (agentId: string, data: Partial<AgentData>) => void
  setAgents: (agents: Record<string, AgentData>) => void
  setLeaderboard: (leaderboard: AgentData[]) => void
  setGitInsights: (insights: GitInsights) => void
  addEvent: (event: any) => void
  setWsConnected: (connected: boolean) => void
  handleWebSocketMessage: (message: any) => void
}

export const useDagStore = create<DagState>((set, get) => ({
  dagData: null,
  taskStates: {},
  agents: {},
  leaderboard: [],
  gitInsights: null,
  gitStats: {
    totalCommits: 0,
    totalLinesAdded: 0,
    totalLinesRemoved: 0,
    fileChangeFrequency: {}
  },
  events: [],
  wsConnected: false,

  setDagData: (data) => set({ dagData: data }),
  
  updateTaskState: (taskId, state) => 
    set((prevState) => ({
      taskStates: {
        ...prevState.taskStates,
        [taskId]: state
      }
    })),
  
  updateTaskStates: (states) => 
    set({ taskStates: states }),
  
  updateAgent: (agentId, data) =>
    set((prevState) => ({
      agents: {
        ...prevState.agents,
        [agentId]: {
          ...prevState.agents[agentId],
          ...data
        }
      }
    })),
  
  setAgents: (agents) => set({ agents }),
  
  setLeaderboard: (leaderboard) => set({ leaderboard }),
  
  setGitInsights: (insights) => set({ gitInsights: insights }),
  
  addEvent: (event) =>
    set((prevState) => ({
      events: [...prevState.events.slice(-99), event]
    })),
  
  setWsConnected: (connected) => set({ wsConnected: connected }),
  
  handleWebSocketMessage: (message) => {
    const state = get()
    
    switch (message.type) {
      case 'init':
        set({
          dagData: message.dagData || null,
          taskStates: message.tasks || {},
          agents: message.agents || {},
          leaderboard: message.leaderboard || [],
          gitInsights: message.gitInsights || null,
          events: message.recentEvents || []
        })
        break
      
      case 'task-update':
        const { task, status, agent, timestamp } = message
        state.updateTaskState(task, { status, agent, timestamp })
        
        // Show toast for task updates
        const statusIcon = status === 'completed' ? '✅' : 
                          status === 'failed' ? '❌' : 
                          status === 'started' ? '🚀' : '⚡'
        toast(`${statusIcon} Task ${task} ${status}`, {
          description: agent ? `Agent: ${agent}` : undefined,
          duration: 3000,
        })
        break
      
      case 'task-claimed':
        const { task: claimedTask, agent: claimAgent } = message
        state.updateTaskState(claimedTask, { 
          status: 'started', 
          agent: claimAgent,
          timestamp: message.timestamp
        })
        
        // Show toast for task claimed
        toast(`🎯 Task ${claimedTask} claimed`, {
          description: `Agent: ${claimAgent}`,
          duration: 3000,
        })
        break
      
      case 'git-commit':
        // Handle git commit updates
        state.addEvent({ type: 'commit', ...message })
        
        // Show toast for git commits
        const commitMessage = message.message || message.data?.message || 'No message'
        const commitAgent = message.agent || 'Unknown'
        const filesChanged = message.data?.files?.length || message.filesChanged || 0
        
        toast(`🔀 New commit by ${commitAgent}`, {
          description: `${commitMessage.substring(0, 50)}${commitMessage.length > 50 ? '...' : ''} (${filesChanged} files)`,
          duration: 4000,
        })
        break
      
      case 'stats-update':
        if (message.leaderboard) {
          state.setLeaderboard(message.leaderboard)
        }
        if (message.gitInsights) {
          state.setGitInsights(message.gitInsights)
        }
        break
      
      case 'event':
        state.addEvent(message)
        break
      
      default:
        console.log('Unknown message type:', message.type)
    }
  }
}))

export default useDagStore
export const useStore = useDagStore