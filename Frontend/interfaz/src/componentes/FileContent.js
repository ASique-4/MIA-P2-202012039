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
    const requestData = {
      comando: fileContent
    };
    console.log(requestData);
  
    const options = {
      method: 'POST',
      body: JSON.stringify(requestData) // Convertir a cadena JSON
    };
  
    fetch('http://localhost:8080/ejecutar-comando', options)
      .then(response => response.json())
      .then(response => console.log(response.mensaje))
      .catch(err => console.error(err));
  };

  return (
    <div>
      <label htmlFor="file-upload" className='custom-file-upload'>
        <i className="fa fa-cloud-upload"></i> Subir Archivo
      </label>
      <input type='file' id='file-upload' onChange={e => handleFileChosen(e.target.files[0])} />
      <input type='button' value='Ejecutar' id='Procesar' className='Procesar' onClick={e => handleButtonClick()} />
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