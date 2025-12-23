import "./Landing.css"
import { Link } from "react-router-dom"

const Landing = () => {
  return (
    <div className="landing">
      <div className="get-started">
        <button><Link to="/auth">Get Started</Link></button>
        <button><Link to="/chats">Chats</Link></button>
      </div>
    </div>
  )
}

export default Landing