const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"

function getToken(){
    return localStorage.getItem("token")
}

async function request(path: string, options: RequestInit = {}){
    const token = getToken()
    const res = await fetch(`${API_URL}${path}`, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...options.headers,
      },
    })
    if (!res.ok) throw new Error(await res.text())
    return res.json()
}

export const api = {
    register: (email: string, password: string) =>
        request("/auth/register", {method : "POST", body: JSON.stringify({email, password}) }),
    login: (email: string, password: string) =>
      request("/auth/login", { method: "POST", body: JSON.stringify({ email, password }) }),

    getCategories: () => request("/categories"),
    getSummary: () => request("/budgets/summary"),
    getBudgets: () => request("/budgets"),

    createExpense: (data: object) =>
      request("/expenses", { method: "POST", body: JSON.stringify(data) }),

    createBudget: (data: object) =>
      request("/budgets", { method: "POST", body: JSON.stringify(data) }),

    deleteBudget: (id: number) =>
      request(`/budgets/${id}`, { method: "DELETE" }),
}