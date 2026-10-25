"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { api } from "@/lib/api"

type Category = {
  ID: number
  Name: string
}

export default function NewExpensePage() {
  const router = useRouter()
  const [categories, setCategories] = useState<Category[]>([])
  const [categoryId, setCategoryId] = useState("")
  const [amount, setAmount] = useState("")
  const [note, setNote] = useState("")
  const [date, setDate] = useState(new Date().toISOString().split("T")[0])
  const [error, setError] = useState("")

  useEffect(() => {
    if (!localStorage.getItem("token")) {
      router.push("/login")
      return
    }
    api.getCategories().then(setCategories).catch(() => router.push("/login"))
  }, [])

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError("")
    try {
      await api.createExpense({ categoryId: Number(categoryId), amount, note, date })

      router.push("/dashboard")
    } catch {
      setError("Failed to create expense.")
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 to-slate-800">
      <div className="w-full max-w-md bg-white/10 backdrop-blur rounded-2xl shadow-xl p-8 border border-white/10">
        <div className="flex items-center gap-3 mb-6">
          <button
            onClick={() => router.back()}
            className="text-white/50 hover:text-white transition text-sm"
          >
            ← Back
          </button>
          <h1 className="text-xl font-semibold text-white">Add Expense</h1>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1 text-white/70">Category</label>
            <select
              value={categoryId}
              onChange={(e) => setCategoryId(e.target.value)}
              className="w-full bg-slate-800 border border-white/20 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-blue-400"
              required
            >
              <option key="default" value="" className="bg-slate-800">Select a category</option>
              {categories.map((c) => (
                <option key={String(c.ID)} value={c.ID} className="bg-slate-800">{c.Name}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium mb-1 text-white/70">Amount</label>
            <input
              type="number"
              step="0.01"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              className="w-full bg-white/10 border border-white/20 rounded-lg px-3 py-2 text-white placeholder-white/30 focus:outline-none focus:border-blue-400"
              placeholder="0.00"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1 text-white/70">Note</label>
            <input
              type="text"
              value={note}
              onChange={(e) => setNote(e.target.value)}
              className="w-full bg-white/10 border border-white/20 rounded-lg px-3 py-2 text-white placeholder-white/30 focus:outline-none focus:border-blue-400"
              placeholder="Optional"
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1 text-white/70">Date</label>
            <input
              type="date"
              value={date}
              onChange={(e) => setDate(e.target.value)}
              className="w-full bg-white/10 border border-white/20 rounded-lg px-3 py-2 text-white focus:outline-none focus:border-blue-400"
              required
            />
          </div>
          {error && <p className="text-red-400 text-sm">{error}</p>}
          <button
            type="submit"
            className="w-full bg-blue-500 hover:bg-blue-600 text-white py-2 rounded-lg font-medium transition"
          >
            Save Expense
          </button>
        </form>
      </div>
    </div>
  )
}
