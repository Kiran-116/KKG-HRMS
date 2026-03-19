import React, { useEffect, useState, useRef } from 'react';

interface DocumentViewModalProps {
  isOpen: boolean;
  documentId: string;
  fileName: string;
  onClose: () => void;
}

const DocumentViewModal: React.FC<DocumentViewModalProps> = ({
  isOpen,
  documentId,
  fileName,
  onClose,
}) => {
  const [pdfUrl, setPdfUrl] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const blobUrlRef = useRef<string | null>(null);

  useEffect(() => {
    if (isOpen && documentId) {
      loadDocument();
    } else {
      // Cleanup blob URL when modal closes
      if (blobUrlRef.current) {
        URL.revokeObjectURL(blobUrlRef.current);
        blobUrlRef.current = null;
        setPdfUrl(null);
      }
    }

    // Cleanup on unmount
    return () => {
      if (blobUrlRef.current) {
        URL.revokeObjectURL(blobUrlRef.current);
        blobUrlRef.current = null;
      }
    };
  }, [isOpen, documentId]);

  const loadDocument = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Get the API base URL
      const apiBaseUrl = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';
      const token = localStorage.getItem('access_token');
      
      // Fetch the document as blob
      const response = await fetch(`${apiBaseUrl}/documents/${documentId}/download`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Failed to load document');
      }

      const blob = await response.blob();
      const url = URL.createObjectURL(blob);
      blobUrlRef.current = url;
      setPdfUrl(url);
    } catch (err: any) {
      setError(err.message || 'Failed to load document');
    } finally {
      setLoading(false);
    }
  };


  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div className="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        {/* Background overlay */}
        <div
          className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
          onClick={onClose}
        ></div>

        {/* Center modal */}
        <span className="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">
          &#8203;
        </span>

        {/* Modal panel */}
        <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-6xl sm:w-full">
          {/* Header */}
          <div className="bg-white px-4 pt-5 pb-4 sm:px-6 border-b border-gray-200">
            <div className="flex items-center justify-between">
              <h3 className="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                {fileName}
              </h3>
              <button
                type="button"
                onClick={onClose}
                className="text-gray-400 hover:text-gray-500 focus:outline-none"
              >
                <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>

          {/* Content */}
          <div className="bg-white px-4 pt-4 pb-4 sm:px-6 sm:pb-6">
            {loading && (
              <div className="flex items-center justify-center h-96">
                <div className="text-center">
                  <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto"></div>
                  <p className="mt-4 text-sm text-gray-500">Loading document...</p>
                </div>
              </div>
            )}

            {error && (
              <div className="flex items-center justify-center h-96">
                <div className="text-center">
                  <p className="text-red-600">{error}</p>
                </div>
              </div>
            )}

            {pdfUrl && !loading && !error && (
              <div className="w-full h-[600px] border border-gray-200 rounded-lg overflow-hidden">
                <iframe
                  src={pdfUrl}
                  className="w-full h-full"
                  title={fileName}
                />
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default DocumentViewModal;
