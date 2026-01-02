import { NavLink } from "react-router-dom"

export default function Sidebar() {
  const menuClass = ({ isActive }: { isActive: boolean }) =>
    `block px-4 py-3 rounded-lg transition ${
      isActive ? "bg-white/20 font-semibold" : "hover:bg-white/20"
    }`

  return (
    <aside className="w-64 bg-gradient-to-b from-blue-700 to-purple-700 text-white flex flex-col">
      <div className="px-6 py-6 text-xl font-bold border-b border-white/20">
        📚 Digital Library
      </div>

      <nav className="flex-1 px-4 py-6 space-y-2">
        <NavLink to="/books" className={menuClass}>
          📘 Koleksi Buku
        </NavLink>

        <NavLink to="/journals" className={menuClass}>
          📄 Jurnal
        </NavLink>

        <NavLink to="/research" className={menuClass}>
          🔬 Laporan Penelitian
        </NavLink>
      </nav>

      <div className="px-6 py-4 text-sm border-t border-white/20 text-white/80">
        © 2025 Perpustakaan Digital
      </div>
    </aside>
  )
}
