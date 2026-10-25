"use client"

import { createContext, useContext, useEffect, useState } from "react"
import { useWebSocket, WSMessage } from "@/hooks/useWebSocket"

const WSContext = createContext<WSMessage | null>(null)

export function WebSocketProvider({ children }: { children: React.ReactNode }) {
  const [token, setToken] = useState<string | null>(null)

  useEffect(() => {
    setToken(localStorage.getItem("token"))
  }, [])

  const lastMessage = useWebSocket(token)

  return <WSContext.Provider value={lastMessage}>{children}</WSContext.Provider>
}

export function useWSMessage() {
  return useContext(WSContext)
}
