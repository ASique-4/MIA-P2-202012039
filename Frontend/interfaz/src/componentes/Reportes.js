function Reportes() {

    const MostrarDISK = () => {
        const reporte = localStorage.getItem('rutaDISK');
        console.log(reporte);
        const reporteDISK = document.getElementById('reportes');
        
        if (reporte) {
          const extension = reporte.split('.').pop().toLowerCase();
          console.log(extension + ': jpg');
          
          if (extension === 'pdf') {
            // Mostrar archivo PDF usando un elemento <embed>
            reporteDISK.innerHTML = `<embed src="${reporte}" type="application/pdf" width="100%" height="600px" />`;
          } else if (extension === 'jpg' || extension === 'jpeg' || extension === 'png') {
            // Mostrar imagen usando un elemento <img>
            reporteDISK.innerHTML = `<img src="${reporte}" alt="Reporte DISK" />`;
          } else {
            // Extensión de archivo no compatible
            reporteDISK.innerHTML = 'Archivo no compatible';
          }
        } else {
          // No hay valor en localStorage
          reporteDISK.innerHTML = 'No se encontró el reporte';
        }
      };
      


  return (
    <div>
      <center>
        <div className="w-full max-w-xs" style={{ margin: "100px" }}>
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
            <button className="Procesar" onClick={e => MostrarDISK()}>DISK</button>
          </div>
          <h3
            className="block text-gray-700 text-xl font-bold mb-2"
            style={{ margin: "25px 0px 25px 0px" }}
          >
            Fase2
          </h3>
          <div className=" content-center">
            <button className="Procesar">TREE</button>
            <button className="Procesar">FILE</button>
            <button className="Procesar">SB</button>
          </div>
          <div id="reportes"></div>
        </div>
      </center>
    </div>
  );
}

export default Reportes;
