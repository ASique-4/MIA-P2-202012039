import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Login from './componentes/Login';
import Reportes from './componentes/Reportes';
import ComplexNavbar from './componentes/NavBar';
import MostrarReportes from './componentes/mostrarReporte';


export default function Index() {
  // Limpiar el local storage
  // Obtener el usuario del local storage
  const user = localStorage.getItem('user');
  console.log(user);

  return (
    <BrowserRouter>
    <ComplexNavbar username={user}/>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/reportes" element={<Reportes />} />
        <Route path="/login" element={<Login />} />
        <Route path="/mostrarReporte" element={<MostrarReportes />} />
      </Routes>
    </BrowserRouter>
  );
}

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(<Index />);
