import "./Register.css"
import { useState } from "react";
import { Register as GoRegister } from "../../../../wailsjs/go/main/App";
import { passwordHash } from "../../../utils/passwordHash";
import { validatePassword } from "../../../utils/validatePassword";
import { validateEmail } from "../../../utils/validateEmail";

const Register = ({changeAuth}) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confPassword, setConfPassword] = useState("");
  const [errMsg, setErrMsg] = useState("");

  const handleChangeAuth = () => {
    changeAuth(true);
  }

  const handleSubmit = async (e) => {
    e.preventDefault();
    setErrMsg("");

    if (password !== confPassword || !validateEmail(email) || !validatePassword(password)) {
      setErrMsg("Invalid data")
      return;
    }

    const hashedPassword = passwordHash(password);

    try {
      const res = await GoRegister(email, hashedPassword);

      console.log("Server says: ", res);
    } catch (err) {
      console.log(err);
      setErrMsg(err);
    }
  }

  return (
    <div className='register'>
      <div className="register-head">
        <p>Create an account.</p>
      </div>
      <div className="register-input">
        {errMsg && <div style={{backgroundColor: "#1A1A1A"}}>
          <br/>
          <p style={{color: "red"}}>{errMsg}</p>
        </div>}

        <input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="Email"/>
        <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="Password"/>

        <input type="password" value={confPassword} onChange={(e) => setConfPassword(e.target.value)} placeholder="Confirm Password"
        style={{outline: confPassword && confPassword !== password ? "1px solid red" : ""}}/>

        <button type="submit" onClick={handleSubmit}>Continue</button>
        <p onClick={handleChangeAuth}>Already have an account?</p>        
      </div>
    </div>
  )
}

export default Register