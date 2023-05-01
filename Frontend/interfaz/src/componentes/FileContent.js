import React, { useState } from 'react';
import PopupComponent  from './Popup';

function FileContent() {
  const [fileContent, setFileContent] = useState('');
  const [showPopup, setShowPopup] = useState(false);
  const [popupTitle, setPopupTitle] = useState('');

  const handleOpenPopup = () => {
    setShowPopup(true);
  };

  const handleAccept = () => {
    console.log('Popup accepted');
    handleConfirm(true);
  };

  const handleReject = () => {
    console.log('Popup rejected');
    handleConfirm(false);
  };



  const handleFileRead = (e) => {
    const content = e.target.result;
    setFileContent(content);
  }

  const handleFileChosen = (file) => {
    console.log("handleFileChosen")
    let fileReader = new FileReader();
    console.log(file);
    fileReader.onloadend = handleFileRead;
    fileReader.readAsText(file);
  }


/**
 * Esta función envía una solicitud POST a un servidor con una carga JSON y maneja la respuesta.
 * @param respuesta - `respuesta` es una variable que contiene el valor de la respuesta del usuario a
 * un mensaje o diálogo de confirmación. Se pasa como argumento a la función `handleConfirm`. Luego, la
 * función usa este valor para crear un objeto de datos de solicitud que se envía a un servidor usando
 * la API `fetch`
 */
  const handleConfirm = (respuesta) => {
   
    const requestData = {
      Aceptar: respuesta
    };

    const options = {
      method: 'POST',
      body: JSON.stringify(requestData) // Convertir a cadena JSON
    };

    setPopupTitle('¿Está seguro de que desea eliminar el disco?');

    const salida = document.getElementById('salida');
    fetch('http://localhost:3030/confirmar', options)
      .then(response => response.json())
      .then(response => {
        console.log(response);
        salida.innerText += "=======> " + response.accion + " <=======\n";
        salida.innerText += response.mensaje + '\n';
      })

  }

  const obtenerRutaReporte = (requestData) => {
    // Se obtiene la ruta del reporte
    // El comando es de la forma: rep >id=061A >Path=/home/user/reports/reporte1.jpg >name=mbr
    // Se separa el comando por el caracter '>' para obtener cada uno de los parametros
    const comando = requestData.comando;
    const comandoSeparado = comando.split('>');
    // Se obtiene el parametro 'Path' sin importar el orden en que se envie
    const parametroPath = comandoSeparado.find(parametro => parametro.includes('Path'));
    // Se obtiene el valor del parametro 'Path'
    const valorPath = parametroPath.split('=')[1];
    return valorPath;
  }

  const imprimirConsola = (accion, mensaje) => {
    const salida = document.getElementById('salida');
    salida.innerText += "=======> " + accion + " <=======\n";
    salida.innerText += mensaje + '\n';
  }


  
  // Enviar archivo al servidor
  const handleButtonClick = () => {
    const fileContent = document.getElementById('file-content').innerText;
    console.log(fileContent);
    // Se separa el contenido del archivo por saltos de linea
    const lineas = fileContent.split('\n');
    console.log(lineas);
  
    // Para cada linea se envia al servidor
    lineas.forEach((linea, index) => {
      setTimeout(() => {
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
          .then(response => {
            console.log(response);
            console.log('prueba');
            if (response.accion === 'pause') {
              setPopupTitle('Ejecución en pausada');
              handleOpenPopup();
            } else if (response.accion === 'Eliminando disco...') {
              setPopupTitle('¿Está seguro de que desea eliminar el disco?');
              handleOpenPopup();
            } else if (response.mensaje === 'Reporte DISK generado con éxito') {
              console.log('rutaDISK');
              localStorage.setItem('rutaDISK', obtenerRutaReporte(requestData));
              imprimirConsola(response.accion, response.mensaje);
            } else if (response.mensaje === 'Reporte TREE generado con éxito') {
              localStorage.setItem('rutaTREE', obtenerRutaReporte(requestData));
              imprimirConsola(response.accion, response.mensaje);
            } else if (response.mensaje === 'Reporte FILE generado con éxito') {
              localStorage.setItem('rutaSB', obtenerRutaReporte(requestData));
              imprimirConsola(response.accion, response.mensaje);
            } else if (response.accion === 'Reporte SB generado con éxito') {
              localStorage.setItem('rutaSB', obtenerRutaReporte(requestData));
              imprimirConsola(response.accion, response.mensaje);
            } else {
              imprimirConsola(response.accion, response.mensaje);
            }
          })
          .catch(err => console.error(err));
      }, index * 2000); // Retraso de 2 segundos (2000 ms) entre cada iteración
    });
  };
  

  return (
    <div>
      <PopupComponent
        openPopup={showPopup}
        onAccept={handleAccept}
        onReject={handleReject}
        title={popupTitle}
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