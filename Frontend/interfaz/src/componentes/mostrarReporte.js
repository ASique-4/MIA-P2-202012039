import React from 'react';

function MostrarReportes(reporte) {
    // Cambiamos de ventana
    const reporteDISK = document.getElementById('reportes');

    const mostrar = (ruta, reporteName) => {
        const reporte = `./img/${ruta}`;
        console.log(reporte);
        const reporteDISK = document.getElementById('reportes');
        
        if (reporte) {
          const extension = reporte.split('.').pop().toLowerCase().trim();
          
          
          if (extension === 'pdf') {
            // Mostrar archivo PDF usando un elemento <embed>
            reporteDISK.innerHTML = `<embed src="${reporte}" type="application/pdf" width="100%" height="600px" />`;
          } else if (extension === 'jpg' || extension === 'jpeg' || extension === 'png') {
            // Mostrar imagen usando un elemento <img>
            reporteDISK.innerHTML = ` <h3 class="block text-gray-700 text-xl font-bold mb-2">${ruta}</h3> <img src="${reporte}" alt="Reporte ${reporteName}" />`;
          } else {
            // Extensión de archivo no compatible
            reporteDISK.innerHTML = 'Archivo no compatible';
          }
        } else {
          // No hay valor en localStorage
          reporteDISK.innerHTML = 'No se encontró el reporte';
        }
    }
    
    if (reporte === 'DISK') {
        mostrar(localStorage.getItem('rutaDISK'), 'DISK');
    } else if (reporte === 'TREE') {
        mostrar('rutaTREE', 'TREE');
    } else if (reporte === 'FILE') {
        mostrar('rutaFILE', 'FILE');
    } else if (reporte === 'SB') {
        mostrar('rutaSB', 'SB');
    } else {
        console.log(reporte)
        reporteDISK.innerHTML = 'Archivo no compatible';
    }
    

    return (
        <div>
            <div id="reportes">
            </div>
        </div>
    );
}

export default MostrarReportes;