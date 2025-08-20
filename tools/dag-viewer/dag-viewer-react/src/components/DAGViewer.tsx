import React, { useEffect, useRef, useState } from 'react'
import cytoscape from 'cytoscape'
import dagre from 'cytoscape-dagre'
import { Card, CardContent, CardHeader, CardTitle } from './ui/card'
import { Badge } from './ui/badge'
import { Button } from './ui/button'
import { HoverCard, HoverCardContent, HoverCardTrigger } from './ui/hover-card'
import { Maximize2, Minimize2, ZoomIn, ZoomOut, Move } from 'lucide-react'

// Register the dagre layout for perfect DAG positioning
cytoscape.use(dagre)

interface DAGViewerProps {
  dagData: any
  taskStates: Record<string, any>
  className?: string
}

// Colorblind-friendly palette using patterns and high contrast
// Based on Paul Tol's colorblind-safe palette and IBM Design Language
const stateColors: Record<string, string> = {
  pending: '#F0F0F0',      // Very light gray (neutral)
  started: '#FFE66D',      // Bright yellow (high visibility)
  in_progress: '#4B7CC8',  // Strong blue
  failed: '#1A1A1A',       // Near black (high contrast)
  blocked: '#FF6B35',      // Vermillion orange (distinct)
  completed: '#2E8B57'     // Sea green (distinct from red/orange)
}

const stateBorderColors: Record<string, string> = {
  pending: '#999999',      // Medium gray
  started: '#FFD700',      // Gold (darker yellow)
  in_progress: '#003F88',  // Dark blue
  failed: '#000000',       // Black
  blocked: '#CC4125',      // Dark vermillion
  completed: '#1B5E3F'     // Dark sea green
}

// High contrast patterns for wave grouping
// Using distinct shades and patterns
const waveColors = [
  '#000000', // Black
  '#4472C4', // Blue
  '#70AD47', // Green  
  '#FFC000', // Amber
  '#7030A0', // Purple
  '#ED7D31', // Orange
  '#A5A5A5', // Gray
  '#255E91', // Navy
  '#44723C', // Forest
  '#CC9900', // Dark amber
]

export const DAGViewer: React.FC<DAGViewerProps> = ({ dagData, taskStates, className }) => {
  const containerRef = useRef<HTMLDivElement>(null)
  const cyRef = useRef<cytoscape.Core | null>(null)
  const [selectedNode, setSelectedNode] = useState<any>(null)
  const [hoverCardOpen, setHoverCardOpen] = useState(false)
  const [isFullscreen, setIsFullscreen] = useState(false)
  const [clickPosition, setClickPosition] = useState({ x: 0, y: 0 })

  useEffect(() => {
    if (!containerRef.current || !dagData) return

    // Convert DAG data to Cytoscape format
    const elements: any[] = []

    // Add nodes from topo_order
    if (dagData.topo_order) {
      dagData.topo_order.forEach((taskId: string) => {
        const state = taskStates[taskId]?.status || 'pending'
        // Get task metadata if available
        const taskMeta = dagData.tasks?.find((t: any) => t.id === taskId) || 
                        dagData.task_metadata?.[taskId]
        
        // Find which wave this task belongs to
        let waveNumber = null
        if (dagData.waves) {
          const wave = dagData.waves.find((w: any) => w.tasks?.includes(taskId))
          waveNumber = wave?.waveNumber
        }
        
        // Get feature info if available
        let featureInfo = null
        if (taskMeta?.feature_id && dagData.features) {
          featureInfo = dagData.features.find((f: any) => f.id === taskMeta.feature_id)
        }
        
        elements.push({
          data: {
            id: taskId,
            label: taskId, // Keep short ID for node display
            title: taskMeta?.title || '',
            description: taskMeta?.description || '',
            category: taskMeta?.category || '',
            duration: taskMeta?.duration || null,
            durationUnits: taskMeta?.durationUnits || '',
            skillsRequired: taskMeta?.skillsRequired || [],
            feature_id: taskMeta?.feature_id || '',
            feature_title: featureInfo?.title || '',
            interfaces_produced: taskMeta?.interfaces_produced || [],
            interfaces_consumed: taskMeta?.interfaces_consumed || [],
            waveNumber,
            taskState: state
          }
        })
      })
    }

    // Add edges from either reduced_edges_sample or dependencies
    const edges = dagData.reduced_edges_sample || 
                  dagData.dependencies?.map((d: any) => [d.from, d.to]) || []
    
    edges.forEach(([source, target]: [string, string]) => {
      elements.push({
        data: {
          id: `${source}-${target}`,
          source,
          target
        }
      })
    })

    // Initialize Cytoscape
    const cy = cytoscape({
      container: containerRef.current,
      elements,
      style: [
        {
          selector: 'node',
          style: {
            'label': 'data(label)',
            'text-valign': 'center',
            'text-halign': 'center',
            'background-color': (ele: any) => {
              const state = ele.data('taskState')
              return stateColors[state] || '#94a3b8'
            },
            'border-color': (ele: any) => {
              const state = ele.data('taskState')
              // Use wave color for border if available and task is pending
              const waveNum = ele.data('waveNumber')
              if (waveNum && state === 'pending') {
                return waveColors[(waveNum - 1) % waveColors.length]
              }
              return stateBorderColors[state] || '#64748b'
            },
            'border-width': (ele: any) => {
              const waveNum = ele.data('waveNumber')
              return waveNum ? 4 : 3
            },
            'color': (ele: any) => {
              // Use contrasting text colors based on background
              const state = ele.data('taskState')
              if (state === 'pending' || state === 'started') {
                return '#000000' // Black text on light backgrounds
              }
              return '#FFFFFF' // White text on dark backgrounds
            },
            'font-size': '12px',
            'font-weight': 'bold',
            'width': 85,
            'height': 45,
            'text-wrap': 'wrap',
            'text-max-width': '75px',
            'overlay-color': (ele: any) => {
              const waveNum = ele.data('waveNumber')
              if (waveNum) {
                return waveColors[(waveNum - 1) % waveColors.length]
              }
              return 'transparent'
            },
            'overlay-opacity': 0.15,
            'overlay-padding': 3
          }
        },
        {
          selector: 'edge',
          style: {
            'width': 2.5,
            'line-color': '#404040',      // Dark gray for high contrast
            'target-arrow-color': '#404040',
            'target-arrow-shape': 'triangle',
            'curve-style': 'bezier',
            'arrow-scale': 1.2             // Larger arrows for visibility
          }
        },
        {
          selector: '.selected',
          style: {
            'border-width': 6,
            'border-color': '#FF1744',    // Bright red for high visibility
            'overlay-color': '#FF1744',
            'overlay-padding': 6,
            'overlay-opacity': 0.25,
            'z-index': 999                // Bring selected node to front
          }
        }
      ],
      layout: {
        name: 'dagre',
        rankDir: 'TB', // Top to bottom layout
        animate: false,
        fit: true,
        padding: 50,
        nodeSep: 50, // Space between nodes in the same rank
        rankSep: 100, // Space between ranks
        edgeSep: 10, // Space between edges
        ranker: 'network-simplex', // Algorithm for rank assignment
        nodeDimensionsIncludeLabels: true
      },
      wheelSensitivity: 0.3,
      minZoom: 0.1,
      maxZoom: 3
    })

    cyRef.current = cy

    // Handle node clicks for hover card
    cy.on('tap', 'node', (evt) => {
      const node = evt.target
      const renderedPosition = node.renderedPosition()
      const containerRect = containerRef.current?.getBoundingClientRect()
      
      if (containerRect) {
        setClickPosition({
          x: renderedPosition.x + containerRect.left,
          y: renderedPosition.y + containerRect.top
        })
      }
      
      cy.nodes().removeClass('selected')
      node.addClass('selected')
      const nodeData = node.data()
      setSelectedNode({
        ...nodeData,
        taskData: taskStates[nodeData.id]
      })
      setHoverCardOpen(true)
    })

    // Handle background clicks
    cy.on('tap', (evt) => {
      if (evt.target === cy) {
        cy.nodes().removeClass('selected')
        setSelectedNode(null)
        setHoverCardOpen(false)
      }
    })

    return () => {
      cy.destroy()
    }
  }, [dagData])

  // Update node states when taskStates change
  useEffect(() => {
    if (!cyRef.current) return

    Object.entries(taskStates).forEach(([taskId, stateData]) => {
      const node = cyRef.current?.getElementById(taskId)
      if (node && node.length > 0) {
        const state = stateData.status || 'pending'
        node.data('taskState', state)
        
        // Trigger style update
        node.style({
          'background-color': stateColors[state] || '#94a3b8',
          'border-color': stateBorderColors[state] || '#64748b'
        })
      }
    })
  }, [taskStates])

  const getStateCounts = () => {
    const counts: Record<string, number> = {
      pending: 0,
      started: 0,
      in_progress: 0,
      failed: 0,
      blocked: 0,
      completed: 0
    }

    Object.values(taskStates).forEach((stateData: any) => {
      const state = stateData.status || 'pending'
      if (counts[state] !== undefined) {
        counts[state]++
      }
    })

    return counts
  }

  const getGraphStats = () => {
    const nodeCount = dagData?.topo_order?.length || 0
    const edgeCount = dagData?.reduced_edges_sample?.length || 0
    const waveCount = dagData?.metrics?.longestPath || 0
    const density = dagData?.metrics?.edgeDensity || 0
    
    return {
      nodes: nodeCount,
      edges: edgeCount,
      waves: waveCount,
      density: (density * 100).toFixed(1)
    }
  }

  const handleResetZoom = () => {
    if (cyRef.current) {
      cyRef.current.fit()
    }
  }

  const handleZoomIn = () => {
    if (cyRef.current) {
      const currentZoom = cyRef.current.zoom()
      cyRef.current.zoom({
        level: currentZoom * 1.2,
        renderedPosition: { x: containerRef.current!.clientWidth / 2, y: containerRef.current!.clientHeight / 2 }
      })
    }
  }

  const handleZoomOut = () => {
    if (cyRef.current) {
      const currentZoom = cyRef.current.zoom()
      cyRef.current.zoom({
        level: currentZoom * 0.8,
        renderedPosition: { x: containerRef.current!.clientWidth / 2, y: containerRef.current!.clientHeight / 2 }
      })
    }
  }

  const handleToggleFullscreen = () => {
    setIsFullscreen(!isFullscreen)
    setTimeout(() => {
      if (cyRef.current) {
        cyRef.current.resize()
        cyRef.current.fit()
      }
    }, 100)
  }

  const stateCounts = getStateCounts()
  const graphStats = getGraphStats()

  const graphContent = (
    <>
      {/* Color Legend */}
      <Card className="mb-2">
        <CardContent className="p-3">
          <div className="flex flex-wrap gap-3 items-center">
            <span className="text-xs font-semibold text-muted-foreground">Legend:</span>
            {Object.entries(stateColors).map(([state, color]) => (
              <div key={state} className="flex items-center gap-1">
                <div 
                  className="w-4 h-4 rounded border-2" 
                  style={{ 
                    backgroundColor: color,
                    borderColor: stateBorderColors[state]
                  }}
                />
                <span className="text-xs capitalize">{state.replace('_', ' ')}</span>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
      
      <Card className="overflow-hidden">
        <CardHeader className="pb-3">
          <div className="flex justify-between items-center">
            <div className="flex gap-4">
              <div className="text-sm text-muted-foreground">
                <span className="font-semibold">{graphStats.nodes}</span> nodes
              </div>
              <div className="text-sm text-muted-foreground">
                <span className="font-semibold">{graphStats.edges}</span> edges
              </div>
              <div className="text-sm text-muted-foreground">
                <span className="font-semibold">{graphStats.waves}</span> waves
              </div>
              <div className="text-sm text-muted-foreground">
                <span className="font-semibold">{graphStats.density}%</span> density
              </div>
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="icon"
                onClick={handleZoomOut}
                title="Zoom Out"
              >
                <ZoomOut className="h-4 w-4" />
              </Button>
              <Button
                variant="outline"
                size="icon"
                onClick={handleZoomIn}
                title="Zoom In"
              >
                <ZoomIn className="h-4 w-4" />
              </Button>
              <Button
                variant="outline"
                size="icon"
                onClick={handleResetZoom}
                title="Reset Zoom"
              >
                <Move className="h-4 w-4" />
              </Button>
              <Button
                variant="outline"
                size="icon"
                onClick={handleToggleFullscreen}
                title={isFullscreen ? "Exit Fullscreen" : "Fullscreen"}
              >
                {isFullscreen ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent className="p-0">
          <div 
            ref={containerRef} 
            className={`w-full bg-white dark:bg-gray-900 ${
              isFullscreen ? 'h-[calc(100vh-200px)]' : 'h-[700px]'
            }`}
            style={{
              background: 'linear-gradient(45deg, #FFFFFF 25%, #F5F5F5 25%, #F5F5F5 50%, #FFFFFF 50%, #FFFFFF 75%, #F5F5F5 75%, #F5F5F5)',
              backgroundSize: '20px 20px'
            }}
          />
        </CardContent>
      </Card>

      {/* Hover Card for Task Details */}
      {selectedNode && (
        <div 
          style={{ 
            position: 'fixed', 
            left: clickPosition.x, 
            top: clickPosition.y,
            zIndex: 9999
          }}
        >
          <HoverCard open={hoverCardOpen} onOpenChange={setHoverCardOpen}>
            <HoverCardTrigger asChild>
              <div />
            </HoverCardTrigger>
            <HoverCardContent className="w-96 max-h-[600px] overflow-y-auto" align="center" side="right">
              <div className="space-y-3">
                {/* Header with ID and Title */}
                <div>
                  <h4 className="text-sm font-semibold">{selectedNode.id}</h4>
                  {selectedNode.title && (
                    <p className="text-base font-medium mt-1">{selectedNode.title}</p>
                  )}
                  {selectedNode.description && (
                    <p className="text-sm text-muted-foreground mt-2">{selectedNode.description}</p>
                  )}
                  <div className="flex gap-2 mt-2">
                    <Badge 
                      style={{ 
                        backgroundColor: stateColors[selectedNode.taskState],
                        borderColor: stateBorderColors[selectedNode.taskState]
                      }}
                    >
                      {selectedNode.taskState?.replace('_', ' ').toUpperCase()}
                    </Badge>
                    {selectedNode.category && (
                      <Badge variant="outline">
                        {selectedNode.category}
                      </Badge>
                    )}
                  </div>
                </div>

                {/* Execution Info */}
                {(selectedNode.taskData?.agent || selectedNode.taskData?.startTime) && (
                  <div className="space-y-2 pt-2 border-t">
                    <p className="text-xs font-semibold">Execution</p>
                    {selectedNode.taskData?.agent && (
                      <div className="flex items-center gap-2">
                        <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
                        <span className="text-sm">Agent: {selectedNode.taskData.agent}</span>
                      </div>
                    )}
                    {selectedNode.taskData?.startTime && (
                      <div className="text-xs text-muted-foreground">
                        Started: {new Date(selectedNode.taskData.startTime * 1000).toLocaleTimeString()}
                      </div>
                    )}
                    {selectedNode.taskData?.tokensSpent !== undefined && selectedNode.taskData.tokensSpent > 0 && (
                      <div className="text-xs text-muted-foreground">
                        Tokens Used: {selectedNode.taskData.tokensSpent.toLocaleString()}
                      </div>
                    )}
                  </div>
                )}

                {/* Duration Estimates */}
                {selectedNode.duration && (
                  <div className="space-y-2 pt-2 border-t">
                    <p className="text-xs font-semibold">Duration Estimates</p>
                    {typeof selectedNode.duration === 'object' ? (
                      <>
                        <div className="text-xs text-muted-foreground">
                          Optimistic: {selectedNode.duration.optimistic} {selectedNode.durationUnits || 'hours'}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          Most Likely: {selectedNode.duration.mostLikely} {selectedNode.durationUnits || 'hours'}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          Pessimistic: {selectedNode.duration.pessimistic} {selectedNode.durationUnits || 'hours'}
                        </div>
                      </>
                    ) : (
                      <div className="text-xs text-muted-foreground">
                        Estimated: {selectedNode.duration} {selectedNode.durationUnits || 'hours'}
                      </div>
                    )}
                  </div>
                )}

                {/* Skills Required */}
                {selectedNode.skillsRequired && selectedNode.skillsRequired.length > 0 && (
                  <div className="space-y-2 pt-2 border-t">
                    <p className="text-xs font-semibold">Skills Required</p>
                    <div className="flex flex-wrap gap-1">
                      {selectedNode.skillsRequired.map((skill: string) => (
                        <Badge key={skill} variant="secondary" className="text-xs">
                          {skill}
                        </Badge>
                      ))}
                    </div>
                  </div>
                )}

                {/* Wave Info */}
                {selectedNode.waveNumber && (
                  <div className="space-y-1 pt-2 border-t">
                    <p className="text-xs font-semibold">Wave</p>
                    <div className="text-xs text-muted-foreground">
                      Wave {selectedNode.waveNumber}
                    </div>
                  </div>
                )}

                {/* Feature Info */}
                {selectedNode.feature_id && (
                  <div className="space-y-1 pt-2 border-t">
                    <p className="text-xs font-semibold">Feature</p>
                    <div className="text-xs text-muted-foreground">
                      {selectedNode.feature_id}
                      {selectedNode.feature_title && (
                        <div className="mt-1">{selectedNode.feature_title}</div>
                      )}
                    </div>
                  </div>
                )}

                {/* Interfaces */}
                {(selectedNode.interfaces_produced?.length > 0 || selectedNode.interfaces_consumed?.length > 0) && (
                  <div className="space-y-2 pt-2 border-t">
                    <p className="text-xs font-semibold">Interfaces</p>
                    {selectedNode.interfaces_produced?.length > 0 && (
                      <div>
                        <span className="text-xs text-muted-foreground">Produces:</span>
                        <div className="flex flex-wrap gap-1 mt-1">
                          {selectedNode.interfaces_produced.map((iface: string) => (
                            <Badge key={iface} variant="outline" className="text-xs">
                              {iface}
                            </Badge>
                          ))}
                        </div>
                      </div>
                    )}
                    {selectedNode.interfaces_consumed?.length > 0 && (
                      <div>
                        <span className="text-xs text-muted-foreground">Consumes:</span>
                        <div className="flex flex-wrap gap-1 mt-1">
                          {selectedNode.interfaces_consumed.map((iface: string) => (
                            <Badge key={iface} variant="outline" className="text-xs">
                              {iface}
                            </Badge>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>
                )}

                {/* Dependencies */}
                <div className="pt-2 border-t">
                  <p className="text-xs font-semibold mb-1">Dependencies</p>
                  <div className="space-y-1">
                    {/* Check both formats */}
                    {(dagData?.reduced_edges_sample || dagData?.dependencies)?.filter((edge: any) => {
                      const [_, to] = Array.isArray(edge) ? edge : [edge.from, edge.to]
                      return to === selectedNode.id
                    }).map((edge: any) => {
                      const [from] = Array.isArray(edge) ? edge : [edge.from]
                      const depTask = dagData?.tasks?.find((t: any) => t.id === from)
                      return (
                        <div key={from} className="text-xs text-muted-foreground">
                          ← {from} {depTask?.title && <span className="text-xs">({depTask.title})</span>}
                        </div>
                      )
                    })}
                    {(dagData?.reduced_edges_sample || dagData?.dependencies)?.filter((edge: any) => {
                      const [from, to] = Array.isArray(edge) ? edge : [edge.from, edge.to]
                      return from === selectedNode.id
                    }).map((edge: any) => {
                      const [, to] = Array.isArray(edge) ? edge : [edge.from, edge.to]
                      const depTask = dagData?.tasks?.find((t: any) => t.id === to)
                      return (
                        <div key={to} className="text-xs text-muted-foreground">
                          → {to} {depTask?.title && <span className="text-xs">({depTask.title})</span>}
                        </div>
                      )
                    })}
                  </div>
                </div>
              </div>
            </HoverCardContent>
          </HoverCard>
        </div>
      )}
    </>
  )

  return (
    <div className={className}>
      {/* Only show status badges if there are any tasks with status */}
      {Object.values(stateCounts).some(count => count > 0) && (
        <Card className="mb-4">
          <CardContent className="p-4">
            <div className="flex gap-2 flex-wrap">
              {Object.entries(stateCounts).map(([state, count]) => (
                count > 0 && (
                  <Badge key={state} variant="secondary" className="gap-2">
                    <div 
                      className="w-2 h-2 rounded-full" 
                      style={{ backgroundColor: stateColors[state] }}
                    />
                    <span className="capitalize">{state.replace('_', ' ')}: {count}</span>
                  </Badge>
                )
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {isFullscreen ? (
        <div className="fixed inset-0 z-50 bg-background p-4">
          {graphContent}
        </div>
      ) : (
        graphContent
      )}
    </div>
  )
}