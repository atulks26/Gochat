
import './App.css';
import { BrowserRouter, MemoryRouter, Routes, Route, useNavigate, Link, Navigate } from "react-router-dom"
import Landing from './components/landing/Landing';
import Auth from './components/auth/Auth';
import Chats from './components/chats/Chats';
import { useAuth, AuthProvider } from './context/userContext';

const ProtectedRoute = ({ children }) => {
    const { user } = useAuth();
    if (!user) {
        return <Navigate to="/auth" />;
    }
    
    return children;
};

function App() {
    return (
        <AuthProvider>
            <BrowserRouter>
                <button style={{border: "none", backgroundColor: "green", padding: "1rem"}}><Link to="/">Home</Link></button>
                <Routes>
                    <Route path='/' element={<Landing />}/>
                    <Route path='/auth' element={<Auth />}/>
                    <Route path='/chats' element={<Chats />}/>
                </Routes>
            </BrowserRouter>
        </AuthProvider>
        
    )
}

export default App
