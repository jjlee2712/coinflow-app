"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { api } from "@/lib/api"

type Budget = {
  ID: number
  CategoryID: number
  Amount: string
}

type Category = {
  ID: number
  Name: string
}

export default function BudgetsPage() {
  const router = useRouter()
  const [budgets, setBudgets] = useState<Budget[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [categoryId, setCategoryId] = useState("")
  const [amount, setAmount] = useState("")
  const [error, setError] = useState("")

  useEffect(() => {
    if (!localStorage.getItem("token")) {
      router.push("/login")
      return
    }
    loadData()
  }, [])

  async function loadData() {
    const [b, c] = await Promise.all([api.getBudgets(), api.getCategories()])
    setBudgets(b || [])
    setCategories(c || [])
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    setError("")
    try {
      await api.createBudget({ categoryId: Number(categoryId), amount })
      setAmount("")
      setCategoryId("")
      loadData()
    } catch {
      setError("Failed to create budget. This category may already have one.")
    }
  }

  async function handleDelete(id: number) {
    await api.deleteBudget(id)
    loadData()
  }

  function getCategoryName(id: number) {
    return categories.find((c) => c.ID === id)?.Name ?? `Category ${id}`
  }



  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800">
      {/* Navbar */}
      <nav className="border-b border-white/10 px-6 py-4">
        <div className="max-w-3xl mx-auto flex justify-between items-center">
          <h1 className="text-xl font-bold text-white">💰 Coinflow</h1>
          <a href="/dashboard" className="text-white/50 hover:text-white text-sm transition">
            ← Dashboard
          </a>
        </div>
      </nav>

      <div className="max-w-3xl mx-auto px-6 py-8">
        <h2 className="text-white text-2xl font-bold mb-6">Manage Budgets</h2>

        {/* Create form */}
        <div className="bg-white/10 backdrop-blur rounded-2xl p-6 mb-6 border border-white/10">
          <h3 className="text-white font-semibold mb-4">Set a Budget</h3>
          <form onSubmit={handleCreate} className="flex gap-3 flex-wrap">
            <select
              value={categoryId}
              onChange={(e) => setCategoryId(e.target.value)}
              className="bg-slate-800 border border-white/20 rounded-lg px-3 py-2 text-white flex-1 min-w-[150px] focus:outline-none focus:border-blue-400"
              required
            >
              <option value="">Select category</option>
              {categories.map((c) => (
                <option key={String(c.ID)} value={c.ID} className="bg-slate-800">
                  {c.Name}
                </option>
              ))}
            </select>
            <input
              type="number"
              step="0.01"
              placeholder="Monthly limit"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              className="bg-white/10 border border-white/20 rounded-lg px-3 py-2 text-white placeholder-white/30 flex-1 min-w-[120px] focus:outline-none focus:border-blue-400"
              required
            />
            <button
              type="submit"
              className="bg-blue-500 hover:bg-blue-600 text-white px-5 py-2 rounded-lg text-sm font-medium transition"
            >
              Add
            </button>
          </form>
          {error && <p className="text-red-400 text-sm mt-2">{error}</p>}
        </div>

        {/* Budget list */}
        <div className="space-y-3">
          {budgets.length === 0 ? (
            <div className="bg-white/10 rounded-2xl p-8 text-center">
              <p className="text-white/50">No budgets yet. Create one above.</p>
            </div>
          ) : (
            budgets.map((b) => (
              <div
                key={b.ID}
                className="bg-white/10 backdrop-blur rounded-2xl p-5 flex justify-between items-center border border-white/10 hover:bg-white/15 transition"
              >
                <div>
                  <p className="text-white font-medium">{getCategoryName(b.CategoryID)}</p>
                  <p className="text-white/50 text-sm">RM {b.Amount} / month</p>
                </div>
                <button
                  onClick={() => handleDelete(b.ID)}
                  className="text-red-400 hover:text-red-300 text-sm transition"
                >
                  Delete
                </button>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  )
}
