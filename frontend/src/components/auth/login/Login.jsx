import "./Login.css"
import { useRef, useState } from "react";
import { Login as GoLogin } from "../../../../wailsjs/go/main/App";
import { useAuth } from "../../../context/userContext";

const Login = ({changeAuth}) => {
  const usernameRef = useRef(null);
  const passwordRef = useRef(null); 
  const [errMsg, setErrMsg] = useState("");
  const {login} = useAuth();
  
  const handleChangeAuth = () => {
    changeAuth(false);
  }

  const handleSubmit = async (e) => {
    e.preventDefault();
    setErrMsg("");

    const username = usernameRef.current.value;
    const password = passwordRef.current.value;

    if (!username || !password) return;

    try {
      const user = await GoLogin(username, password);
      console.log("Server says: ", user);

      login(user);
      //navigate to chats after login
    } catch (err) {
      console.error(err);
      setErrMsg(err.message);
    }
  }

  return (
    <div className='login'>
      <div className="login-head">
        <p>Sign back in.</p>
      </div>
      
      <div className="login-input">
        {errMsg && <div style={{backgroundColor: "#1A1A1A"}}>
          <br/>
          <p style={{color: "red"}}>{errMsg}</p>
        </div>}

        <input ref={usernameRef} placeholder="Username or Email"/>
        <input type="password" ref={passwordRef} placeholder="Password"/>
        <button onClick={handleSubmit}>Continue</button>
        <p onClick={handleChangeAuth}>Not registered yet?</p>
      </div>
    </div>
  )
}

export default Login