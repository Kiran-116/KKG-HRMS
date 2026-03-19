import React, { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { documentService, DocumentWithUser } from '../../services/documentService';
import ConfirmModal from '../../components/ConfirmModal';
import DocumentViewModal from '../../components/DocumentViewModal';

const AdminDocumentListPage: React.FC = () => {
  const [documents, setDocuments] = useState<DocumentWithUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const limit = 10;
  const [confirmModal, setConfirmModal] = useState<{
    isOpen: boolean;
    documentId: string | null;
  }>({
    isOpen: false,
    documentId: null,
  });
  const [viewModal, setViewModal] = useState<{
    isOpen: boolean;
    documentId: string | null;
    fileName: string;
  }>({
    isOpen: false,
    documentId: null,
    fileName: '',
  });

  useEffect(() => {
    loadDocuments();
  }, [page]);

  const loadDocuments = async () => {
    try {
      setLoading(true);
      const data = await documentService.getAllDocuments(page, limit);
      setDocuments(data.documents as DocumentWithUser[]);
      setTotal(data.total);
    } catch (error) {
      console.error('Failed to load documents:', error);
      toast.error('Failed to load documents');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteClick = (id: string) => {
    setConfirmModal({
      isOpen: true,
      documentId: id,
    });
  };

  const handleConfirmDelete = async () => {
    if (!confirmModal.documentId) return;

    try {
      await documentService.delete(confirmModal.documentId);
      toast.success('Document deleted successfully!');
      loadDocuments();
    } catch (error: any) {
      toast.error(error.response?.data?.message || 'Failed to delete document');
    } finally {
      setConfirmModal({ isOpen: false, documentId: null });
    }
  };

  const handleCancelDelete = () => {
    setConfirmModal({ isOpen: false, documentId: null });
  };

  const handleView = (id: string, fileName: string) => {
    setViewModal({
      isOpen: true,
      documentId: id,
      fileName: fileName,
    });
  };

  const handleCloseView = () => {
    setViewModal({
      isOpen: false,
      documentId: null,
      fileName: '',
    });
  };

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
  };

  if (loading) {
    return <div className="text-center py-12">Loading documents...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold text-gray-900">All Documents</h1>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">File Name</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Employee Name</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Size</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Uploaded</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {documents?.map((doc) => (
              <tr key={doc.id}>
                <td className="px-6 py-4 text-sm text-gray-900">{doc.file_name}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{doc.user_name || '-'}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{doc.user_email || '-'}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{doc.document_type}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatFileSize(doc.file_size)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {new Date(doc.uploaded_at).toLocaleDateString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <div className="flex space-x-4">
                    <button
                      onClick={() => handleView(doc.id, doc.file_name)}
                      className="text-indigo-600 hover:text-indigo-900"
                    >
                      View
                    </button>
                    <button
                      onClick={() => handleDeleteClick(doc.id)}
                      className="text-red-600 hover:text-red-900"
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {documents?.length === 0 && (
        <div className="text-center py-12 text-gray-500">No documents found</div>
      )}

      <div className="flex justify-between items-center">
        <p className="text-sm text-gray-600">Showing {documents.length} of {total} documents</p>
        <div className="flex space-x-2">
          <button
            onClick={() => setPage(p => Math.max(1, p - 1))}
            disabled={page === 1}
            className="px-4 py-2 border rounded-md disabled:opacity-50"
          >
            Previous
          </button>
          <button
            onClick={() => setPage(p => p + 1)}
            disabled={documents.length < limit}
            className="px-4 py-2 border rounded-md disabled:opacity-50"
          >
            Next
          </button>
        </div>
      </div>

      <ConfirmModal
        isOpen={confirmModal.isOpen}
        title="Delete Document"
        message="Are you sure you want to delete this document? This action cannot be undone."
        confirmText="Delete"
        cancelText="Cancel"
        type="danger"
        onConfirm={handleConfirmDelete}
        onCancel={handleCancelDelete}
      />

      {viewModal.documentId && (
        <DocumentViewModal
          isOpen={viewModal.isOpen}
          documentId={viewModal.documentId}
          fileName={viewModal.fileName}
          onClose={handleCloseView}
        />
      )}
    </div>
  );
};

export default AdminDocumentListPage;
