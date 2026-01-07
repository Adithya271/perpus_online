import { useEffect, useState } from "react"
import Sidebar from "../components/layout/Sidebar"
import { api } from "../services/api"

type Recommendation = {
  id_buku: number
  judul: string
  penulis: string
  tahun: number
  status: string
}

export default function Recommendations() {
  const [recommendations, setRecommendations] = useState<Recommendation[]>([])
  const [loading, setLoading] = useState(true)

  const userId = 1

  useEffect(() => {
    loadRecommendations()
  }, [])

  const loadRecommendations = async () => {
    setLoading(true)
    try {
      const res = await api.get(`/recommendations/${userId}/by-search`)
      setRecommendations(res.data?.recommendations || [])
    } catch (err) {
      console.error("Error fetching recommendations:", err)
      setRecommendations([])
    } finally {
      setLoading(false)
    }
  }

  const handleBorrow = async (bookId: number) => {
    try {
      await api.post("/loans", {
        id_user: userId,
        id_buku: bookId,
      })
      alert("Berhasil meminjam buku!")
      loadRecommendations()
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any) {
      alert(err.response?.data?.error || "Gagal meminjam buku")
    }
  }

  return (
    <div className="flex min-h-screen bg-gray-100">
      <Sidebar />

      <main className="flex-1 p-8">
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-gray-800">
             Rekomendasi Buku
          </h1>
          <p className="text-gray-500">Berdasarkan riwayat pencarian Anda</p>
        </div>

        {loading ? (
          <div className="flex flex-col items-center justify-center py-24">
            <div className="animate-spin h-12 w-12 rounded-full border-b-2 border-blue-600" />
            <p className="mt-4 text-gray-500 font-medium">
              Menganalisis minat Anda...
            </p>
          </div>
        ) : recommendations.length === 0 ? (
          <div className="bg-white rounded-xl p-12 text-center shadow border">
            <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg
                className="w-8 h-8 text-blue-600"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
            </div>
            <h3 className="text-lg font-bold text-gray-800">
              Belum Ada Rekomendasi
            </h3>
            <p className="text-gray-500 mt-2 max-w-md mx-auto">
              Silakan cari beberapa buku di halaman koleksi agar sistem dapat
              memberikan rekomendasi yang sesuai.
            </p>
            <a
              href="/"
              className="inline-block mt-6 px-6 py-3 bg-blue-600 text-white rounded-lg font-semibold hover:bg-blue-700"
            >
              Kembali ke Koleksi Buku
            </a>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
            {recommendations.map((book, index) => (
              <div
                key={book.id_buku}
                className="bg-white rounded-xl shadow-sm border hover:shadow-md transition flex flex-col"
              >
                <div className="p-5 flex-grow">
                  <div className="flex justify-between mb-3">
                    <span className="text-xs font-bold bg-blue-100 text-blue-700 px-2 py-1 rounded">
                      REKOMENDASI #{index + 1}
                    </span>
                  </div>

                  <h3 className="font-bold text-gray-900 text-lg leading-tight mb-1 uppercase">
                    {book.judul}
                  </h3>
                  <p className="text-sm text-gray-500 mb-4">
                    {book.penulis} * {book.tahun}
                  </p>
                </div>

                <div className="p-4 border-t bg-gray-50 flex justify-between items-center">
                  <span
                    className={`text-xs font-bold ${
                      book.status === "tersedia"
                        ? "text-green-600"
                        : "text-red-600"
                    }`}
                  >
                    * {book.status.toUpperCase()}
                  </span>

                  <button
                    onClick={() => handleBorrow(book.id_buku)}
                    disabled={book.status !== "tersedia"}
                    className={`px-4 py-2 rounded-lg text-sm font-semibold transition ${
                      book.status === "tersedia"
                        ? "bg-blue-600 text-white hover:bg-blue-700"
                        : "bg-gray-200 text-gray-400 cursor-not-allowed"
                    }`}
                  >
                    Pinjam
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>
    </div>
  )
}
