import { useEffect, useRef, useState, useCallback } from 'react'

interface WebSocketOptions {
  url: string
  onMessage?: (data: any) => void
  onOpen?: () => void
  onClose?: () => void
  onError?: (error: Event) => void
  reconnectInterval?: number
  maxReconnectAttempts?: number
}

export const useWebSocket = (options: string | WebSocketOptions) => {
  const config: WebSocketOptions = typeof options === 'string' 
    ? { url: options }
    : options;
    
  const {
    url,
    onMessage,
    onOpen,
    onClose,
    onError,
    reconnectInterval = 5000,
    maxReconnectAttempts = 5
  } = config;
  const [isConnected, setIsConnected] = useState(false)
  const [reconnectAttempts, setReconnectAttempts] = useState(0)
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const shouldReconnectRef = useRef(true)

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN || !shouldReconnectRef.current) {
      return
    }

    if (reconnectAttempts >= maxReconnectAttempts) {
      console.log(`Max reconnection attempts (${maxReconnectAttempts}) reached. Stopping.`)
      return
    }

    try {
      const ws = new WebSocket(url)
      
      ws.onopen = () => {
        console.log('WebSocket connected')
        setIsConnected(true)
        setReconnectAttempts(0)
        onOpen?.()
        
        if (reconnectTimeoutRef.current) {
          clearTimeout(reconnectTimeoutRef.current)
          reconnectTimeoutRef.current = null
        }
      }
      
      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          onMessage?.(data)
        } catch (error) {
          // Silently handle parse errors
        }
      }
      
      ws.onerror = (error) => {
        // Silently handle errors to avoid console spam
        onError?.(error)
      }
      
      ws.onclose = () => {
        setIsConnected(false)
        onClose?.()
        wsRef.current = null
        
        // Attempt to reconnect with exponential backoff
        if (shouldReconnectRef.current && reconnectAttempts < maxReconnectAttempts) {
          const delay = Math.min(reconnectInterval * Math.pow(1.5, reconnectAttempts), 30000)
          console.log(`WebSocket disconnected. Reconnecting in ${(delay/1000).toFixed(1)}s (attempt ${reconnectAttempts + 1}/${maxReconnectAttempts})`)
          
          reconnectTimeoutRef.current = setTimeout(() => {
            setReconnectAttempts(prev => prev + 1)
            connect()
          }, delay)
        }
      }
      
      wsRef.current = ws
    } catch (error) {
      setIsConnected(false)
    }
  }, [url, onMessage, onOpen, onClose, onError, reconnectInterval, maxReconnectAttempts, reconnectAttempts])

  const disconnect = useCallback(() => {
    shouldReconnectRef.current = false
    
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }
    
    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }
    
    setIsConnected(false)
    setReconnectAttempts(0)
  }, [])

  const sendMessage = useCallback((data: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(data))
      return true
    }
    return false
  }, [])

  useEffect(() => {
    shouldReconnectRef.current = true
    connect()
    
    return () => {
      disconnect()
    }
  }, []) // Remove dependencies to avoid reconnection loops

  return {
    isConnected,
    sendMessage,
    connect,
    disconnect,
    reconnectAttempts
  }
}