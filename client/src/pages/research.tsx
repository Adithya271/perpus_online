import Sidebar from "../components/layout/Sidebar"

export default function Research() {
  return (
    <div className="flex min-h-screen bg-gray-100">
      <Sidebar />
      <main className="flex-1 p-8">
        <h1 className="text-2xl font-bold">🔬 Laporan Penelitian</h1>
        <p className="mt-4 text-gray-600">
          Halaman penelitian (dalam pengembangan)
        </p>
      </main>
    </div>
  )
}
