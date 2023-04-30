import React, { useState } from 'react';
import PopupComponent  from './Popup';

function FileContent() {
  const [fileContent, setFileContent] = useState('');
  const [showPopup, setShowPopup] = useState(false);

  const handleOpenPopup = () => {
    setShowPopup(true);
  };

  const handleAccept = () => {
    console.log('Popup accepted');
  };

  const handleReject = () => {
    console.log('Popup rejected');
  };



  const handleFileRead = (e) => {
    const content = e.target.result;
    setFileContent(content);
  }

  const handleFileChosen = (file) => {
    let fileReader = new FileReader();
    fileReader.onloadend = handleFileRead;
    fileReader.readAsText(file);
  }

  function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }


  
  // Enviar archivo al servidor
  const handleButtonClick = () => {
    const fileContent = document.getElementById('file-content').innerText;
    console.log(fileContent);
    // Se separa el contenido del archivo por saltos de linea
    const lineas = fileContent.split('\n');
    console.log(lineas);

    // Para cada linea se envia al servidor
    lineas.forEach(linea => {
      const salida = document.getElementById('salida');
      console.log(linea);
      // Si es un comentario
      if (linea.startsWith('#')) {
        salida.innerText += linea + '\n';
        return;
      }

      // Si es una linea vacia
      if (linea === '') {
        return;
      }

      const requestData = {
        comando: linea
      };
      console.log(requestData);
    
      const options = {
        method: 'POST',
        body: JSON.stringify(requestData) // Convertir a cadena JSON
      };
  
    
      fetch('http://localhost:8080/ejecutar-comando', options)
        .then(response => response.json())
        .then(response => response.accion === "pause" ? handleOpenPopup() : (salida.innerText += '======   ' + response.accion + '   ======\n' + response.mensaje + '\n'))
        .then(sleep(1000))
        .catch(err => console.error(err));
    }
    );

    
  };

  return (
    <div>
      <PopupComponent
        openPopup={showPopup}
        onAccept={handleAccept}
        onReject={handleReject}
        title='El servidor se encuentra en pausa'
      />
      <label htmlFor="file-upload" className='custom-file-upload'>
        <i className="fa fa-cloud-upload"></i> Subir Archivo
      </label>
      <input type='file' id='file-upload' onChange={e => handleFileChosen(e.target.files[0])} />
      <input type='button' value='Ejecutar' id='Procesar' className='Procesar' onClick={e => handleButtonClick()} />
      <input type="button" value='Limpiar' id='Limpiar' className='Limpiar' onClick={e => document.getElementById('salida').innerText = ''} />
      <pre 
        id='file-content'
        style={{
          overflowY: 'scroll',
          height: '300px',
          scrollbarWidth: 'thin',
          scrollbarColor: 'transparent transparent',
          msOverflowStyle: 'none',
          WebkitOverflowScrolling: 'touch',
        }}
      >
  {fileContent}
</pre>
<div>
      <pre
        id='salida'
        style={{
          overflowY: 'scroll',
          height: '200px',
          scrollbarWidth: 'thin',
          scrollbarColor: 'transparent transparent',
          msOverflowStyle: 'none',
          WebkitOverflowScrolling: 'touch',
        }}
      >
</pre>

    </div>

    </div>
  );
}

export default FileContent;