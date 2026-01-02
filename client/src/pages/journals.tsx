import Sidebar from "../components/layout/Sidebar"

export default function Journals() {
  return (
    <div className="flex min-h-screen bg-gray-100">
      <Sidebar />
      <main className="flex-1 p-8">
        <h1 className="text-2xl font-bold">📄 Jurnal</h1>
        <p className="mt-4 text-gray-600">
          Halaman jurnal (dalam pengembangan)
        </p>
      </main>
    </div>
  )
}
