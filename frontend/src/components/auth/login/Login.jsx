import "./Login.css"
import { useRef, useState } from "react";
import { Login as GoLogin } from "../../../../wailsjs/go/main/App"
import { passwordHash } from "../../../utils/passwordHash";

const Login = ({changeAuth}) => {
  const usernameRef = useRef(null);
  const passwordRef = useRef(null); 
  const [errMsg, setErrMsg] = useState("");
  
  const handleChangeAuth = () => {
    changeAuth(false);
  }

  const handleSubmit = async (e) => {
    e.preventDefault();
    setErrMsg("");

    const username = usernameRef.current.value;
    const password = passwordRef.current.value;

    if (!username || !hashedPassword) return;

    const hashedPassword = passwordHash(password);

    try {
      const res = await GoLogin(username, hashedPassword);

      console.log("Server says: ", res);
    } catch (err) {
      console.error(err);
      setErrMsg(err);
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
        <button type="submit" onClick={handleSubmit}>Continue</button>
        <p onClick={handleChangeAuth}>Not registered yet?</p>
      </div>
    </div>
  )
}

export default Login  