import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import StatusListView from './components/StatusListView';
import StatusDetailView from './components/StatusDetailView';
import StatusForm from './components/StatusForm';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<StatusListView />} />
        <Route path="/detail/:id" element={<StatusDetailView />} />
        <Route path="/create" element={<StatusForm />} />
        <Route path="/edit/:id" element={<StatusForm />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
