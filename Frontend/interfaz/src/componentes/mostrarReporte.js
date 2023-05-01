import React from 'react';

function MostrarReportes(reporte) {
  const reporteDiv = document.getElementById('reportes');

  const mostrar = (base, reporteName) => {
    const imageUrl = `data:image/jpeg;base64,${base}`;
    reporteDiv.innerHTML = `
      <h3 class="block text-gray-700 text-xl font-bold mb-2" style={{ margin: "25px 0px 25px 0px" }}>
        ${reporteName}
      </h3>
      <div>
        <img src="${imageUrl}" alt="${reporteName}">
      </div>
    `;
  }
  

  if (reporte === 'DISK') {
    mostrar(localStorage.getItem('baseDISK'), 'DISK');
  } else if (reporte === 'TREE') {
    mostrar(localStorage.getItem('baseTREE'), 'TREE');
  } else if (reporte === 'FILE') {
    mostrar(localStorage.getItem('baseFILE'), 'FILE');
  } else if (reporte === 'SB') {
    mostrar(localStorage.getItem('baseSB'), 'SB');
  } else {
    console.log(reporte);
    reporteDiv.innerHTML = 'Archivo no compatible';
  }

  return (
    <div>
      <div id="reportes"></div>
    </div>
  );
}

export default MostrarReportes;
