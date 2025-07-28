// DOM要素の取得
const uploadArea = document.getElementById('uploadArea');
const fileInput = document.getElementById('fileInput');
const compressBtn = document.getElementById('compressBtn');
const loader = document.getElementById('loader');
const result = document.getElementById('result');
const error = document.getElementById('error');
const preview = document.getElementById('preview');
const originalSize = document.getElementById('originalSize');
const compressedSize = document.getElementById('compressedSize');
const compressionRatio = document.getElementById('compressionRatio');
const filename = document.getElementById('filename');
const downloadBtn = document.getElementById('downloadBtn');

// ファイル選択UI要素
const uploadIcon = document.getElementById('uploadIcon');
const uploadText = document.getElementById('uploadText');
const fileInfo = document.getElementById('fileInfo');
const fileName = document.getElementById('fileName');
const fileSize = document.getElementById('fileSize');
const changeFileBtn = document.getElementById('changeFileBtn');

let selectedFile = null;

// 定数
const ALLOWED_TYPES = ['image/png', 'image/jpeg', 'image/jpg', 'image/gif'];
const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB

// イベントリスナーの設定
function initializeEventListeners() {
    // ドラッグ&ドロップイベント
    uploadArea.addEventListener('click', () => fileInput.click());
    
    uploadArea.addEventListener('dragover', (e) => {
        e.preventDefault();
        uploadArea.classList.add('dragover');
    });

    uploadArea.addEventListener('dragleave', () => {
        uploadArea.classList.remove('dragover');
    });

    uploadArea.addEventListener('drop', (e) => {
        e.preventDefault();
        uploadArea.classList.remove('dragover');
        const files = e.dataTransfer.files;
        if (files.length > 0) {
            handleFileSelect(files[0]);
        }
    });

    fileInput.addEventListener('change', (e) => {
        if (e.target.files.length > 0) {
            handleFileSelect(e.target.files[0]);
        }
    });

    // 別のファイルを選択ボタンのイベント
    changeFileBtn.addEventListener('click', () => {
        fileInput.click();
    });

    // 圧縮ボタンのイベント
    compressBtn.addEventListener('click', handleCompress);
}

// ファイル選択処理
function handleFileSelect(file) {
    if (!validateFile(file)) {
        return;
    }

    selectedFile = file;
    compressBtn.disabled = false;
    hideError();
    hideResult();

    // UI更新
    updateFileUI(file);

    // プレビュー表示
    showPreview(file);
}

// ファイル検証
function validateFile(file) {
    // ファイル形式チェック
    if (!ALLOWED_TYPES.includes(file.type)) {
        showError('対応していないファイル形式です。PNG, JPG, JPEG, GIF形式のみ対応しています。');
        return false;
    }

    // ファイルサイズチェック
    if (file.size > MAX_FILE_SIZE) {
        showError('ファイルサイズが大きすぎます。10MB以下のファイルを選択してください。');
        return false;
    }

    return true;
}

// ファイルUI更新
function updateFileUI(file) {
    // ファイル情報を表示
    fileName.textContent = file.name;
    fileSize.textContent = formatFileSize(file.size);
    
    // UI状態を更新
    uploadArea.classList.add('has-file');
    uploadIcon.textContent = '✅';
    uploadText.style.display = 'none';
    fileInfo.style.display = 'block';
}

// UIリセット
function resetFileUI() {
    uploadArea.classList.remove('has-file');
    uploadIcon.textContent = '📁';
    uploadText.style.display = 'block';
    fileInfo.style.display = 'none';
    preview.style.display = 'none';
}

// プレビュー表示
function showPreview(file) {
    const reader = new FileReader();
    reader.onload = (e) => {
        preview.src = e.target.result;
        preview.style.display = 'block';
    };
    reader.readAsDataURL(file);
}

// 圧縮処理
async function handleCompress() {
    if (!selectedFile) return;

    // UI状態を更新
    compressBtn.disabled = true;
    showLoader();
    hideError();
    hideResult();

    // FormDataを作成
    const formData = createFormData();

    try {
        const response = await fetch('/api/compress', {
            method: 'POST',
            body: formData
        });

        const data = await response.json();

        if (response.ok && data.success) {
            showResult(data);
        } else {
            showError(data.error || '圧縮に失敗しました。');
        }
    } catch (err) {
        showError('サーバーとの通信に失敗しました。');
    } finally {
        hideLoader();
        compressBtn.disabled = false;
    }
}

// FormData作成
function createFormData() {
    const formData = new FormData();
    formData.append('image', selectedFile);
    
    const quality = document.getElementById('quality').value;
    const width = document.getElementById('width').value;
    const height = document.getElementById('height').value;

    if (quality) formData.append('quality', quality);
    if (width) formData.append('width', width);
    if (height) formData.append('height', height);

    return formData;
}

// 結果表示
function showResult(data) {
    originalSize.textContent = formatFileSize(data.original_size);
    compressedSize.textContent = formatFileSize(data.compressed_size);
    compressionRatio.textContent = data.compression_ratio.toFixed(1) + '%';
    filename.textContent = data.output_file;
    
    downloadBtn.href = data.download_url;
    downloadBtn.style.display = 'inline-block';
    
    result.style.display = 'block';
}

// ローダー表示/非表示
function showLoader() {
    loader.style.display = 'block';
}

function hideLoader() {
    loader.style.display = 'none';
}

// エラー表示/非表示
function showError(message) {
    error.textContent = message;
    error.style.display = 'block';
}

function hideError() {
    error.style.display = 'none';
}

// 結果非表示
function hideResult() {
    result.style.display = 'none';
}

// ファイルサイズフォーマット
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// 初期化
document.addEventListener('DOMContentLoaded', initializeEventListeners); 