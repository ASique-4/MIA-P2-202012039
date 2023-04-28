import './App.css';
import FileContent from './componentes/FileContent';
import ComplexNavbar from './componentes/NavBar';
// Mostrar contenido del archivo


function App() {
  return (
    <div className="App">
    {/* Llamada al componente que muestra el contenido del archivo */}
    <ComplexNavbar username={''} password={''} />
      <FileContent />

    </div>
  );
}

export default App;
