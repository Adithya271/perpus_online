import { useEffect, useState } from "react"
import Sidebar from "../components/layout/Sidebar"
import { api } from "../services/api"
import { useNavigate } from "react-router-dom"


type Book = {
  id_buku: number
  judul: string
  penulis: string
  tahun: number
  status: "tersedia" | "dipinjam" | "habis" | "dikembalikan"
  stok: number
}

type UserLoan = {
  id_buku: number
  status: "dipinjam" | "dikembalikan"
}

export default function Books() {
  const [books, setBooks] = useState<Book[]>([])
  const [userLoans, setUserLoans] = useState<UserLoan[]>([])
  const [loading, setLoading] = useState(true)
  const [keyword, setKeyword] = useState("")
  const [statusFilter, setStatusFilter] = useState("all")
  const navigate = useNavigate()


  const userId = 1 // nanti dari auth

  useEffect(() => {
    loadBooks()
    loadUserLoans()
  }, [])

const loadBooks = async () => {
  setLoading(true)
  try {
    if (keyword.trim() !== "") {
      await saveSearchHistory(keyword) 
    }

    const res = await api.get(`/books?keyword=${keyword}`)
    setBooks(Array.isArray(res.data) ? res.data : [])
  } catch (err) {
    console.error(err)
    setBooks([])
  } finally {
    setLoading(false)
  }
}


  const loadUserLoans = async () => {
    try {
      const res = await api.get(`/loans/user/${userId}`)
      const activeLoans = Array.isArray(res.data)
        ? res.data.filter((l: UserLoan) => l.status === "dipinjam")
        : []
      setUserLoans(activeLoans)
    } catch (err) {
      console.error(err)
      setUserLoans([])
    }
  }

  const handleBorrow = async (id_buku: number) => {
    try {
      await api.post("/loans", { id_user: userId, id_buku })
      alert("Berhasil meminjam buku")
      await loadBooks()
      await loadUserLoans()
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any) {
      alert(err.response?.data?.error || "Gagal meminjam buku")
    }
  }

  const handleReturn = async (id_buku: number) => {
    try {
      await api.post("/loans/return", { id_user: userId, id_buku })
      alert("Buku berhasil dikembalikan")
      await loadBooks()
      await loadUserLoans()
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any) {
      alert(err.response?.data?.error || "Gagal mengembalikan buku")
    }
  }

  const isBorrowedByMe = (id_buku: number) =>
    userLoans.some((l) => l.id_buku === id_buku)

  const filteredBooks = books.filter((b) =>
    statusFilter === "all" ? true : b.status === statusFilter
  )

  const saveSearchHistory = async (keyword: string) => {
    try {
      if (keyword.trim().length < 2) return

      await api.post("/search-history", {
        id_user: userId,
        keyword,
      })
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    } catch (err) {
      console.error("Gagal menyimpan search history")
    }
  }


  return (
    <div className="flex min-h-screen bg-gray-100">
      <Sidebar />

      <main className="flex-1 p-8">
        {/* HEADER */}
        <div className="mb-6 flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <div>
            <h1 className="text-2xl font-bold text-gray-800">
               Koleksi Buku
            </h1>
            <p className="text-gray-500">
              Pencarian dan peminjaman buku perpustakaan
            </p>
          </div>

          {/* TOMBOL REKOMENDASI */}
          <button
            onClick={() => navigate("/recommendations")}
            className="px-5 py-3 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition text-sm font-semibold"
          >
             Rekomendasi Buku
          </button>
        </div>

        {/* SEARCH & FILTER */}
        <div className="bg-white p-6 rounded-xl shadow mb-6">
          <div className="flex flex-col md:flex-row gap-4">
            <input
              type="text"
              placeholder="Cari judul atau penulis..."
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && loadBooks()}
              className="flex-1 px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
            />

            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="px-4 py-3 border rounded-lg"
            >
              <option value="all">Semua Status</option>
              <option value="tersedia">Tersedia</option>
              <option value="habis">Habis</option>
              <option value="dipinjam">Dipinjam</option>
            </select>

            <button
              onClick={loadBooks}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
            >
              Cari
            </button>
          </div>
        </div>

        {/* TABLE */}
        <div className="bg-white rounded-xl shadow overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 border-b">
              <tr>
                <th className="px-6 py-4 text-left">Judul</th>
                <th className="px-6 py-4 text-left">Penulis</th>
                <th className="px-6 py-4 text-center">Tahun</th>
                <th className="px-6 py-4 text-center">Stok</th>
                <th className="px-6 py-4 text-center">Status</th>
                <th className="px-6 py-4 text-center">Aksi</th>
              </tr>
            </thead>

            <tbody className="divide-y">
              {loading ? (
                <tr>
                  <td colSpan={6} className="py-10 text-center text-gray-500">
                    Memuat data...
                  </td>
                </tr>
              ) : filteredBooks.length === 0 ? (
                <tr>
                  <td colSpan={6} className="py-10 text-center text-gray-500">
                    Buku tidak ditemukan
                  </td>
                </tr>
              ) : (
                filteredBooks.map((b) => {
                  const borrowedByMe = isBorrowedByMe(b.id_buku)

                  return (
                    <tr key={b.id_buku} className="hover:bg-blue-50">
                      <td className="px-6 py-4 font-semibold">{b.judul}</td>
                      <td className="px-6 py-4">{b.penulis}</td>
                      <td className="px-6 py-4 text-center">{b.tahun}</td>

                      <td className="px-6 py-4 text-center">
                        <span className="px-3 py-1 rounded-full text-sm bg-gray-100">
                          {b.stok}
                        </span>
                      </td>

                      <td className="px-6 py-4 text-center">
                        <span
                          className={`px-3 py-1 rounded-full text-xs font-semibold ${
                            b.status === "tersedia"
                              ? "bg-green-100 text-green-700"
                              : b.status === "dikembalikan"
                              ? "bg-blue-100 text-blue-700"
                              : "bg-red-100 text-red-700"
                          }`}
                        >
                          {b.status}
                        </span>
                      </td>

                      <td className="px-6 py-4 text-center">
                       
                        {borrowedByMe ? (
                          <button
                            onClick={() => handleReturn(b.id_buku)}
                            className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 text-sm"
                          >
                            Kembalikan
                          </button>
                        ) : b.stok > 0 ? (
                          <button
                            onClick={() => handleBorrow(b.id_buku)}
                            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-sm"
                          >
                            Pinjam
                          </button>
                        ) : (
                          <span className="text-sm text-red-500 font-semibold">
                            Habis
                          </span>
                        )}
                      </td>
                    </tr>
                  )
                })
              )}
            </tbody>
          </table>
        </div>
      </main>
    </div>
  )
}
