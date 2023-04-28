import React, { useState } from 'react';

function FileContent() {
  const [fileContent, setFileContent] = useState('');

  const handleFileRead = (e) => {
    const content = e.target.result;
    setFileContent(content);
  }

  const handleFileChosen = (file) => {
    let fileReader = new FileReader();
    fileReader.onloadend = handleFileRead;
    fileReader.readAsText(file);
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
      console.log(linea);
      const requestData = {
        comando: linea
      };
      console.log(requestData);
    
      const options = {
        method: 'POST',
        body: JSON.stringify(requestData) // Convertir a cadena JSON
      };
  
      const salida = document.getElementById('salida');
    
      fetch('http://52.91.77.62/ejecutar-comando', options)
        .then(response => response.json())
        .then(response => salida.innerText += '======   ' + response.accion + '   ======\n' + response.mensaje + '\n')
        .catch(err => console.error(err));
    }
    );

    
  };

  return (
    <div>
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