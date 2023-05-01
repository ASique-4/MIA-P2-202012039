import React, { useState } from "react";

const Login = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [id, setId] = useState("");

  const handleSubmit = (event) => {
    event.preventDefault();
    const requestData = {
      username: username,
      password: password,
      id: id
    };
    console.log(requestData);

    const options = {
      method: "POST",
      body: JSON.stringify(requestData), // Convertir a cadena JSON
    };

    fetch("http://localhost:8080/login", options)
      .then((response) => response.json())
      .then((response) => {
        if (response.mensaje === "Login correcto") {
          alert("Login exitoso");
          localStorage.setItem("user", username);
        } else {
          alert("Login fallido");
        }
        // Guardamos usuario en el local storage
      })
      .catch((err) => console.error(err));
  };

  return (
    <div>
      <center>
        <div className="w-full max-w-xs" style={{ marginTop: "100px" }}>
          <form
            className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4 glass"
            onSubmit={handleSubmit}
          >
            <div className="mb-4">
              <label
                className="block text-gray-700 text-sm font-bold mb-2"
                htmlFor="is"
              >
                Id
              </label>
              <input
                className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                id="id"
                type="text"
                placeholder="39XX"
                value={id}
                onChange={(e) => setId(e.target.value)}
              />
            </div>
            <div className="mb-4">
              <label
                className="block text-gray-700 text-sm font-bold mb-2"
                htmlFor="username"
              >
                Username
              </label>
              <input
                className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                id="username"
                type="text"
                placeholder="Username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
              />
            </div>
            <div className="mb-6">
              <label
                className="block text-gray-700 text-sm font-bold mb-2"
                htmlFor="password"
              >
                Password
              </label>
              <input
                className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
                id="password"
                type="password"
                placeholder="*********"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
            <div className="flex items-center justify-center">
              <button className="loginEnviar text-xs" type="Submit">
                Iniciar Sesión
              </button>
            </div>
          </form>
        </div>
      </center>
    </div>
  );
};

export default Login;
