import Login from "./login/Login"
import Register from "./register/Register"
import { useState } from "react"
import "./Auth.css"

const Auth = () => {
  const [isRegistered, setIsRegistered] = useState(true);

  return (
    <div className='auth'>
      {isRegistered && <Login changeAuth={setIsRegistered}/>}
      {!isRegistered && <Register changeAuth={setIsRegistered}/>}
    </div>
  )
}

export default Auth