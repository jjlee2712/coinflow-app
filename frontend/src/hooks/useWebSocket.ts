"use client"

import { useEffect, useRef, useState } from "react"

export type WSMessage = {
  type: string
  data: unknown
}

export function useWebSocket(token: string | null) {
  const [lastMessage, setLastMessage] = useState<WSMessage | null>(null)
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    if (!token) return

    const ws = new WebSocket(
      `${process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080"}/ws?token=${token}`
    )
    wsRef.current = ws

    ws.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data)
        setLastMessage(msg)
      } catch {
        // ignore malformed messages
      }
    }

    return () => {
      ws.close()
    }
  }, [token])

  return lastMessage
}
