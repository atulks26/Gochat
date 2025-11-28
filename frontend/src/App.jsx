
import './App.css';
import { BrowserRouter, MemoryRouter, Routes, Route, useNavigate } from "react-router-dom"
import Landing from './components/landing/Landing';
import Auth from './components/auth/Auth';

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path='/' element={<Landing />}/>
                <Route path='/auth' element={<Auth />}/>
            </Routes>
        </BrowserRouter>
    )
}

export default App
