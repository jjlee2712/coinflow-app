"use client"

import { useEffect, useRef, useState } from "react"
import { useRouter } from "next/navigation"

import { api } from "@/lib/api"
import { useWSMessage } from "@/context/WebSocketContext"

type SummaryRow = {
  CategoryID: number
  Categoryname: string
  Budgetamount: string
  Totalspent: string
  Remaining: string
  Percentage: string
}

type Alert = {
  categoryId: number
  categoryName: string
  percentage: number
  severity: string
}

type Toast = {
  id: number
  message: string
  severity: string
}

export default function DashboardPage() {
  const router = useRouter()
  const [summary, setSummary] = useState<SummaryRow[]>([])
  const summaryRef = useRef<SummaryRow[]>([])
  const [toasts, setToasts] = useState<Toast[]>([])
  const lastMessage = useWSMessage()
  const lastHandledMessage = useRef<typeof lastMessage>(null)

  useEffect(() => {
    if (!localStorage.getItem("token")) {
      router.push("/login")
      return
    }
    loadSummary()
  }, [])

  useEffect(() => {
    if (!lastMessage || lastMessage === lastHandledMessage.current) return
    lastHandledMessage.current = lastMessage

    if (lastMessage.type === "summary_update") {
      const rows = lastMessage.data as SummaryRow[]
      setSummary(rows)
      summaryRef.current = rows
    }

    if (lastMessage.type === "budget_alert") {
      const alerts = lastMessage.data as Alert[]
      alerts.forEach((alert) => {
        const id = Date.now() + Math.random()
        const categoryName = alert.categoryName || `Category ${alert.categoryId}`
        const msg =
          alert.severity === "exceeded"
            ? `🚨 Budget exceeded for ${categoryName}! (${alert.percentage.toFixed(1)}% used)`
            : `⚠️ Warning: ${alert.percentage.toFixed(1)}% of budget used for ${categoryName}`

        setToasts((prev) => [...prev, { id, message: msg, severity: alert.severity }])
        setTimeout(() => {
          setToasts((prev) => prev.filter((t) => t.id !== id))
        }, 5000)
      })
    }
  }, [lastMessage])

  async function loadSummary() {
    try {
      const data = await api.getSummary() || []
      setSummary(data)
      summaryRef.current = data
    } catch {
      setSummary([])
    }
  }

  function logout() {
    localStorage.removeItem("token")
    router.push("/login")
  }

  const totalBudget = summary.reduce((s, r) => s + parseFloat(r.Budgetamount || "0"), 0)
  const totalSpent = summary.reduce((s, r) => s + parseFloat(r.Totalspent || "0"), 0)

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800">
      {/* Toasts */}
      <div className="fixed top-4 right-4 space-y-2 z-50">
        {toasts.map((t) => (
          <div
            key={t.id}
            className={`px-4 py-3 rounded-xl shadow-lg text-white text-sm max-w-sm backdrop-blur ${
              t.severity === "exceeded"
                ? "bg-red-600/90"
                : "bg-amber-500/90"
            }`}
          >
            {t.message}
          </div>
        ))}
      </div>

      {/* Navbar */}
      <nav className="border-b border-white/10 px-6 py-4">
        <div className="max-w-5xl mx-auto flex justify-between items-center">
          <h1 className="text-xl font-bold text-white">💰 Coinflow</h1>
          <div className="flex items-center gap-3">
            <a
              href="/expenses/new"
              className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium transition"
            >
              + Add Expense
            </a>
            <a
              href="/budgets"
              className="bg-white/10 hover:bg-white/20 text-white px-4 py-2 rounded-lg text-sm font-medium transition"
            >
              Budgets
            </a>
            <button
              onClick={logout}
              className="text-white/50 hover:text-white text-sm transition"
            >
              Logout
            </button>
          </div>
        </div>
      </nav>

      <div className="max-w-5xl mx-auto px-6 py-8">
        {/* Summary stat cards */}
        <div className="grid grid-cols-3 gap-4 mb-8">
          <div className="bg-white/10 rounded-2xl p-5 backdrop-blur">
            <p className="text-white/60 text-sm mb-1">Total Budget</p>
            <p className="text-white text-2xl font-bold">RM {totalBudget.toFixed(2)}</p>
          </div>
          <div className="bg-white/10 rounded-2xl p-5 backdrop-blur">
            <p className="text-white/60 text-sm mb-1">Total Spent</p>
            <p className="text-white text-2xl font-bold">RM {totalSpent.toFixed(2)}</p>
          </div>
          <div className="bg-white/10 rounded-2xl p-5 backdrop-blur">
            <p className="text-white/60 text-sm mb-1">Remaining</p>
            <p className="text-white text-2xl font-bold">RM {(totalBudget - totalSpent).toFixed(2)}</p>
          </div>
        </div>

        {/* Budget breakdown */}
        <h2 className="text-white/70 text-sm font-medium uppercase tracking-wider mb-4">
          Budget Breakdown
        </h2>

        <div className="space-y-3">
          {summary.length === 0 ? (
            <div className="bg-white/10 rounded-2xl p-8 text-center">
              <p className="text-white/60 mb-3">No budgets set up yet.</p>
              <a
                href="/budgets"
                className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg text-sm font-medium transition"
              >
                Create a Budget
              </a>
            </div>
          ) : (
            summary.map((row) => {
              const pct = parseFloat(row.Percentage) || 0
              const barColor =
                pct >= 100 ? "bg-red-500" : pct >= 80 ? "bg-amber-400" : "bg-emerald-400"
              const badgeColor =
                pct >= 100
                  ? "bg-red-500/20 text-red-400"
                  : pct >= 80
                  ? "bg-amber-500/20 text-amber-400"
                  : "bg-emerald-500/20 text-emerald-400"

              return (
                <div
                  key={row.CategoryID}
                  className="bg-white/10 backdrop-blur rounded-2xl p-5 hover:bg-white/15 transition"
                >
                  <div className="flex justify-between items-start mb-3">
                    <div>
                      <p className="text-white font-semibold">{row.Categoryname}</p>
                      <p className="text-white/50 text-sm">
                        RM {row.Totalspent} spent of RM {row.Budgetamount}
                      </p>
                    </div>
                    <span className={`text-xs font-medium px-2 py-1 rounded-full ${badgeColor}`}>
                      {pct.toFixed(1)}%
                    </span>
                  </div>
                  <div className="w-full bg-white/10 rounded-full h-2">
                    <div
                      className={`h-2 rounded-full transition-all ${barColor}`}
                      style={{ width: `${Math.min(pct, 100)}%` }}
                    />
                  </div>
                  <p className="text-white/40 text-xs mt-2">
                    RM {row.Remaining} remaining
                  </p>
                </div>
              )
            })
          )}
        </div>
      </div>
    </div>
  )
}
