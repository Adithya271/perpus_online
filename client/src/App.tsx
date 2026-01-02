import { BrowserRouter, Routes, Route } from "react-router-dom"
import Books from "./pages/books"
import Journals from "./pages/journals"
import Research from "./pages/research"
import Recommendations from "./pages/recommendations"

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/books" element={<Books />} />
        <Route path="/journals" element={<Journals />} />
        <Route path="/research" element={<Research />} />
        <Route path="/recommendations" element={<Recommendations />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
