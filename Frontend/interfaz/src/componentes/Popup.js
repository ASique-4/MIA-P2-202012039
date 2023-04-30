import React, { useState, useEffect } from 'react';
import Popup from 'reactjs-popup';
import 'reactjs-popup/dist/index.css';

export default function PopupComponent({ openPopup, onAccept, onReject , title}) {
  const [isOpen, setIsOpen] = useState(false);

  const handleOpenPopup = () => {
    setIsOpen(true);
  };

  const handleClosePopup = () => {
    setIsOpen(false);
  };

  const handleAccept = () => {
    handleClosePopup();
    if (onAccept) {
      onAccept();
    }
  };

  const handleReject = () => {
    handleClosePopup();
    if (onReject) {
      onReject();
    }
  };

  // Abre el pop-up cuando se llama a la función openPopup
  useEffect(() => {
    if (openPopup) {
      handleOpenPopup();
    }
  }, [openPopup]);

  return (
    <div className='contenidoPopup'>
      <Popup
        className=' justify-center items-center w-full'
        contentStyle={{
          background: 'rgba( 255, 255, 255, 0.25 )',
          boxShadow: '0 8px 32px 0 #394867',
          backdropFilter: 'blur( 4px )',
          borderRadius: '10px',
          border: '1px solid rgba( 255, 255, 255, 0.18 )',
        }}
        open={isOpen}
        modal
        nested
      >
        {(close) => (
          <div className="modal text-center w-full">
            <div className="content justify-self-center">
              <h4 className='text-xl'>{title}</h4>
            </div>
            <div>
              <button className='Aceptar text-sm' onClick={handleAccept}>Aceptar</button>
              <button className='Rechazar text-sm' onClick={handleReject}>Rechazar</button>
            </div>
          </div>
        )}
      </Popup>
    </div>
  );
}
