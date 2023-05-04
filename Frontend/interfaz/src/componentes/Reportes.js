import React from "react";
import MostrarReportes from "./mostrarReporte";

function Reportes() {

  return (
    <div>
      <center>
        <div className="w-full max-w-xs bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4 glass" style={{ margin: "100px" }}>
          <h1 className="block text-gray-700 text-6xl font-bold mb-2">
            Reportes
          </h1>
          <h3
            className="block text-gray-700 text-xl font-bold mb-2"
            style={{ margin: "30px 0px 25px 0px " }}
          >
            Fase1
          </h3>
          <div className=" content-center">
            <button className="Procesar" onClick={e => MostrarReportes('DISK')}>DISK</button>
          </div>
          <h3
            className="block text-gray-700 text-xl font-bold mb-2"
            style={{ margin: "25px 0px 25px 0px" }}
          >
            Fase2
          </h3>
          <div className=" content-center">
            <button className="Procesar" onClick={e => MostrarReportes('TREE')}>TREE</button>
            <button className="Procesar" onClick={e => MostrarReportes('FILE')}>FILE</button>
            <button className="Procesar" onClick={e => MostrarReportes('SB')}>SB</button>
          </div>
        </div>
        <div id="reportes" style={{margin: "40px"}}>

        </div>
      </center>
    </div>
  );
}

export default Reportes;
