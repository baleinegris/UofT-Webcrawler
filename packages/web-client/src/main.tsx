import { StrictMode } from 'react'
import App from './App'
import ReactDOM from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import './index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <BrowserRouter>
    <StrictMode>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path='*' element={<div className='w-full h-full p-10 bg-gray-300 text-black'>404 Not Found</div>} />
      </Routes>
    </StrictMode>
  </BrowserRouter>
)
 