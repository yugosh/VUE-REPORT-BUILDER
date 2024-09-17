<template>
    <div class="report-builder d-flex flex-row gap-2 justify-content-center">
        <!-- Sidebar Menu -->
        <div class="menu d-flex flex-column gap-2">
            <button class="btn btn-primary" @click="addText">Add Text</button>
            <button class="btn btn-primary mr-2" @click="addImage">Add Image</button>
            <button class="btn btn-primary mr-2" @click="addTable">Add Table</button>

            <button class="btn btn-primary mr-2" @click="addField">Field</button>

            <input type="file" @change="loadDataPayload" accept=".json" />

            <!-- Detail Form -->
            <div class="detail-form" v-if="selectedElement !== null">
                <!-- Form untuk Gambar -->
                <div v-if="elements[selectedElement].type === 'image'">
                    <h3>Edit Image</h3>
                    <label>Image URL:</label>
                    <input type="text" v-model="elements[selectedElement].src" />
                    <label>Width:</label>
                    <input type="number" v-model.number="elements[selectedElement].width" />
                    <label>Height:</label>
                    <input type="number" v-model.number="elements[selectedElement].height" />
                </div>

                <!-- Form untuk Tabel -->
                <div v-if="elements[selectedElement].type === 'table'">
                    <h3>Edit Table Header</h3>
                    <label>Header:</label>
                    <input type="text" v-for="(cell, index) in elements[selectedElement].rows[0]" :key="index"
                        v-model="elements[selectedElement].rows[0][index]" />
                    <h3>Edit Table Placeholders</h3>
                    <label>Placeholders:</label>
                    <input type="text" v-for="(cell, index) in elements[selectedElement].rows[1]" :key="index"
                        v-model="elements[selectedElement].rows[1][index]" />
                    <label>Data Key:</label>
                    <input type="text" v-model="elements[selectedElement].dataKey" />
                    <button @click="addTableHeader">Add Header Row</button>
                    <button @click="addTablePlaceholder">Add Placeholder Row</button>
                    <button @click="addRow">Add Row</button>
                    <button @click="addColumn">Add Column</button>
                    <label>Width:</label>
                    <input type="number" v-model.number="elements[selectedElement].width" />
                    <label>Height:</label>
                    <input type="number" v-model.number="elements[selectedElement].height" />
                </div>

                <!-- Form untuk Teks -->
                <div v-if="elements[selectedElement].type === 'text'">
                    <h3>Edit Text</h3>
                    <label>Content:</label>
                    <input type="text" v-model="elements[selectedElement].content" />
                    <label>Width:</label>
                    <input type="number" v-model.number="elements[selectedElement].width" />
                    <label>Height:</label>
                    <input type="number" v-model.number="elements[selectedElement].height" />

                    <h3>Edit Style</h3>

                    <!-- Bold -->
                    <label>
                        <input type="checkbox" v-model="elements[selectedElement].style.fontWeight" true-value="bold"
                            false-value="normal" />
                        Bold
                    </label>

                    <!-- Italic -->
                    <label>
                        <input type="checkbox" v-model="elements[selectedElement].style.fontStyle" true-value="italic"
                            false-value="normal" />
                        Italic
                    </label>

                    <!-- Underline -->
                    <label>
                        <input type="checkbox" v-model="elements[selectedElement].style.textDecoration"
                            true-value="underline" false-value="none" />
                        Underline
                    </label>

                    <!-- Text Alignment -->
                    <label>Text Alignment:</label>
                    <select v-model="elements[selectedElement].style.textAlign">
                        <option value="left">Left</option>
                        <option value="center">Center</option>
                        <option value="right">Right</option>
                    </select>

                    <!-- Font Size -->
                    <label>Font Size:</label>
                    <input type="number" v-model.number="elements[selectedElement].style.fontSize" />
                </div>
            </div>
        </div>

        <!-- Canvas Area -->
        <div class="canvas-container">
            <div class="canvas" @dragover.prevent @drop="onDrop" ref="canvas">
                <draggable v-model="elements" :itemKey="getItemKey"
                    :options="{ group: 'elements', ghostClass: 'dragging', delay: 0, delayOnTouchOnly: true }">
                    <template v-slot:item="{ element, index }">
                        <div draggable="true" :key="getItemKey(element, index)" @dragstart="onDragStart($event, index)"
                            @click="selectElement(index)" :style="getElementStyle(element)" class="element">
                            <div v-if="element.type === 'text'">
                                <div :style="{
                                        fontWeight: element.style.fontWeight,
                                        fontStyle: element.style.fontStyle,
                                        textDecoration: element.style.textDecoration,
                                        textAlign: element.style.textAlign,
                                        fontSize: element.style.fontSize + 'px'
                                    }"
                                >
                                    {{ parseTextContent(element.content) }}

                                </div>
                            </div>
                            <span v-if="element.type === 'field'">
                                {{ payloadData[element.field] !== undefined ? payloadData[element.field] : `Field
                                ${element.field} not found` }}
                            </span>
                            <img v-if="element.type === 'image'" :src="element.src" alt="Image" />

                            <table v-if="element.type === 'table'" class="table" :style="getTableStyle(element)">
                                <thead>
                                    <tr>
                                        <th v-for="(cell, cellIndex) in element.rows[0]" :key="cellIndex">{{ cell }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr v-for="(row, rowIndex) in getTableData(element.dataKey)" :key="rowIndex">
                                        <td v-for="(placeholder, cellIndex) in element.rows[1]" :key="cellIndex"
                                            :style="getTdStyle()">
                                            {{ row[extractKey(placeholder)] || 'N/A' }}
                                        </td>
                                    </tr>
                                </tbody>
                            </table>

                            <button @click.stop="deleteElement(index)" class="delete-button">X</button>
                        </div>
                    </template>
                </draggable>
            </div>
        </div>

        <div class="menu d-flex flex-column gap-2">
            <button class="btn btn-primary" @click="saveLayout">Save Layout</button>
            <button class="btn btn-primary" @click="generatePDF">Download PDF</button>
            <button class="btn btn-primary" @click="saveAsJSON">Save as JSON</button>
            <input type="file" @change="loadFromJSON" accept=".json" />
        </div>
    </div>
</template>

<script>
import draggable from 'vuedraggable';
import jsPDF from 'jspdf';
import { saveAs } from 'file-saver';
// import html2canvas from 'html2canvas';
import 'jspdf-autotable';

export default {
    name: 'ReportCanvas',
    components: {
        draggable,
    },
    data() {
        return {
            elements: [],   // Struktur elemen di sini atau load dari layout.json yang ada
            selectedElement: null, // Menyimpan indeks elemen yang dipilih
            draggedElementIndex: null,
            canvasWidth: 754, // Ukuran canvas A4 dalam piksel dengan margin 20px
            canvasHeight: 1083,
            margin: 20,
            gridSize: 20, // Ukuran sel grid 20px x 20px
            elementMargin: 4, // Margin antar elemen
            payloadData: null,
        };
    },
    watch: {
        payloadData: {
            handler(newData) {
                console.log('Payload data updated:', newData);
                // this.updateAllTables();
                this.updateCanvas(newData);

                this.elements.forEach((element, index) => {
                    if (element.type === 'variable') {
                        this.elements[index].content = this.getVariableContent(element.key);
                    }
                });
            },
            deep: true
        }
    },
    updated() {
        this.updateCanvas(this.payloadData);
    },
    mounted() {
        // this.generateTableHtml(this.data['elements']);
        this.loadElements();
    },
    methods: {
        loadElements() {
            // Example of loading elements, replace with your actual data source
            // this.elements = [
            //     // Populate with actual elements data
            //     { type: 'text', content: 'Example text', x: 50, y: 50 },
            //     // more elements...
            // ];

            // If the payload data is also ready at this point, update the table
            if (this.payloadData) {
                this.generateTableHtml(this.elements);
            }
        },
        deleteElement(index) {
            if (index >= 0 && index < this.elements.length) {
                this.elements.splice(index, 1);
                if (this.selectedElement === index) {
                    this.selectedElement = null;
                } else if (this.selectedElement > index) {
                    this.selectedElement--;
                }
            }
        },
        addField() {
            let fieldKey = prompt("Enter the key for the field:");
            if (fieldKey) {
                // Hapus $data. jika ada di awal key
                if (fieldKey.startsWith('$data.')) {
                    fieldKey = fieldKey.replace('$data.', '');
                }
                this.elements.push({
                    id: new Date(),
                    type: 'field',
                    field: fieldKey,
                    x: 0,
                    y: 0,
                    width: 200,
                    height: 30
                });
            } else {
                console.error("No key was provided for the field.");
            }
        },
        parseTextContent(text) {
            if (!text) return '';
            return text.replace(/\{\{\s*\$data\.(\w+)\s*\}\}/g, (match, key) => {
                return this.payloadData[key] !== undefined ? this.payloadData[key] : match;
            });
        },
        updateCanvas(data) {
            // Logika untuk mengupdate canvas berdasarkan data yang baru
            // Misalnya, ini bisa jadi metode yang memuat ulang tabel atau grafik di canvas.
            this.payloadData = data;
        },
        updateAllTables() {
            this.elements.forEach(element => {
                if (element.type === 'table') {
                    // Tidak perlu mengubah data elemen, karena kita akan menggunakan payloadData langsung
                }
            });
        },
        loadDataPayload(event) {
            const file = event.target.files[0];
            const reader = new FileReader();
            console.log("payloadDAta", this.payloadData);
            reader.onload = (e) => {
                const data = JSON.parse(e.target.result);
                // this.payloadData = {};
                this.payloadData = data;
                console.log("Loaded data:", this.payloadData);

                this.elements.forEach((element) => {
                    if (element.type === 'table') {
                        // Logika untuk memperbarui tabel jika diperlukan
                    }
                });
            };
            reader.readAsText(file);
        },
        getTableData(dataKey) {
            if (!this.payloadData) return [];
            return dataKey ? this.payloadData[dataKey] || [] : this.payloadData;
        },
        getItemKey(element, index) {
            return element.id || index;
        },
        getElementStyle(element) {
            return {
                position: 'absolute',
                left: `${element.x}px`,
                top: `${element.y}px`,
                width: `${element.width}px`,
                height: `${element.height}px`,
                margin: `${this.elementMargin}px`,
                // border: '1px solid #000', // Tambahkan border ke elemen
                boxSizing: 'border-box',
            };
        },
        getTableStyle(element) {
            return {
                width: `${element.width}px`,
                // height: `${element.height}px`,
                height: element.height ? element.height : 'auto', // Default ke auto jika tidak ada height yang diatur
                // border: '1px dotted #000',
                'border-top': '1px dotted #000',
                'border-left': '1px dotted #000',
                'border-right': '1px dotted #000',
                borderCollapse: 'collapse',
                tableLayout: 'fixed'

            };
        },
        getTdStyle() {
            return {
                border: '1px dotted #000',
                padding: '4px',
                textAlign: 'center',
                verticalAlign: 'top',
                boxSizing: 'border-box',
                minHeight: '20px'
            };
        },
        onDragStart(event, index) {
            console.log('Drag started:', event);
            this.draggedElementIndex = index;

            // Offset position untuk memperbaiki posisi elemen saat di-drag
            const rect = event.target.getBoundingClientRect();
            this.dragOffsetX = event.clientX - rect.left;
            this.dragOffsetY = event.clientY - rect.top;
        },
        onDrop(event) {
            const canvasRect = this.$refs.canvas.getBoundingClientRect();
            let x = event.clientX - canvasRect.left - this.dragOffsetX;
            let y = event.clientY - canvasRect.top - this.dragOffsetY;

            // Snap posisi ke grid
            x = Math.round(x / this.gridSize) * this.gridSize;
            y = Math.round(y / this.gridSize) * this.gridSize;

            this.updateElementPosition(x, y);
        },
        updateElementPosition(x, y) {
            if (this.draggedElementIndex !== null) {
                const element = this.elements[this.draggedElementIndex];

                const elementWidth = element.width + this.elementMargin * 2;
                const elementHeight = element.height + this.elementMargin * 2;

                x = Math.max(this.margin, Math.min(x, this.canvasWidth - elementWidth));
                y = Math.max(this.margin, Math.min(y, this.canvasHeight - elementHeight));

                // Cek apakah elemen baru akan menimpa elemen lain
                const isOverlapping = this.elements.some((el, index) => {
                    if (index === this.draggedElementIndex) return false;
                    return (
                        x < el.x + el.width + this.elementMargin &&
                        x + element.width + this.elementMargin > el.x &&
                        y < el.y + el.height + this.elementMargin &&
                        y + element.height + this.elementMargin > el.y
                    );
                });

                if (!isOverlapping) {
                    this.elements[this.draggedElementIndex].x = x;
                    this.elements[this.draggedElementIndex].y = y;
                }

                this.draggedElementIndex = null;
            }
        },
        saveLayout() {
            console.log('Layout saved:', JSON.stringify(this.elements));
        },
        generatePDF() {
            this.elements.sort((a, b) => a.y - b.y);

            // console.log("elementss", this.elements);

            let htmlContent = `
<html>
  <head>
    <style>
      @font-face {
        font-family: 'CustomFont';
        src: url('https://example.com/path-to-font.ttf');
      }
      body {
        font-family: 'CustomFont';
        margin: 0;
        padding: 0;
      }
      @page {
    size: A4;
    margin: 0;
    }
      .page {
        width: 754px;
        height: 1080px;
        position: relative;
        page-break-after: always;
      }
      .canvas-element {
        position: absolute;
      }
      table {
        width: 100%;
        border-collapse: collapse;
      }
      td, th {
        border: 1px solid black;
        padding: 8px;
      }
      .header {
        position: absolute;
        top: 20px;
        left: 0;
        right: 0;
      }
      .footer {
        position: absolute;
        bottom: 20px;
        left: 0;
        right: 0;
      }
      .table-container {
        position: relative;
        padding-bottom: 60px;
      }
    </style>
  </head>
  <body>`;

            let currentPage = 1;
            let yOffset = 0;
            const pageHeight = 1080;
            let pageContent = '';
            let headerContent = '';
            let footerContent = '';
            let tableFound = false; // Flag untuk menandai apakah tabel sudah ditemukan

            const addElementToPage = (element) => {
                let elementHtml = '';
                const adjustedY = tableFound ? yOffset : element.y - (currentPage - 1) * pageHeight;

                if (element.type === 'text' || element.type === 'field') {
                    const fontWeight = element.style.fontWeight || 'normal';
    const fontStyle = element.style.fontStyle || 'normal';
    const textDecoration = element.style.textDecoration || 'none';
    const textAlign = element.style.textAlign || 'left';
    const fontSize = element.style.fontSize ? `${element.style.fontSize}px` : '14px';
    
                    elementHtml = `
        <div class="canvas-element" style="left: ${element.x}px; 
                top: ${adjustedY}px; 
                font-size: ${fontSize}; 
                font-weight: ${fontWeight}; 
                font-style: ${fontStyle}; 
                text-decoration: ${textDecoration}; 
                text-align: ${textAlign};">
          ${element.content || `{${element.field}}`}
        </div>`;
                } else if (element.type === 'image') {
                    elementHtml = `
        <div class="canvas-element" style="left: ${element.x}px; top: ${adjustedY}px;">
          <img src="${element.src}" width="${element.width}px" height="${element.height}px" />
        </div>`;
                } else if (element.type === 'table') {
                    elementHtml = generateTableHtml(element, adjustedY);
                    tableFound = true; // Tandai bahwa tabel telah ditemukan
                }

                if (element.isHeader) {
                    headerContent += elementHtml;
                } else if (element.isFooter) {
                    footerContent += elementHtml;
                } else {
                    pageContent += elementHtml;
                }

                // Hanya update yOffset jika elemen ditemukan setelah tabel
                if (tableFound) {
                    yOffset = adjustedY + (element.height || 0) + 20; // Update yOffset setelah menambahkan elemen
                }
            };

            const generateTableHtml = (element, adjustedY) => {
                let tableHtml = '';
                const rowsPerPage = 15;
                const totalRows = this.payloadData[element.dataKey]?.length || 0;
                const totalPages = Math.ceil(totalRows / rowsPerPage);
                const rowHeight = 20; // Tinggi rata-rata per baris tabel
                let tableHeight = 0;

                for (let page = 0; page < totalPages; page++) {
                    if (page > 0) {
                        createNewPage();
                        adjustedY = element.y; // Reset adjustedY untuk halaman baru
                    }

                    tableHtml = `
        <div class="canvas-element table-container" style="left: ${element.x}px; top: ${adjustedY}px; width: ${element.width}px;">
          <table>
            <thead>
              <tr>${element.rows[0].map(header => `<th>${header}</th>`).join('')}</tr>
            </thead>
            <tbody>`;

                    const startRow = page * rowsPerPage;
                    const endRow = Math.min(startRow + rowsPerPage, totalRows);

                    if (this.payloadData && this.payloadData[element.dataKey]) {
                        for (let i = startRow; i < endRow; i++) {
                            const row = this.payloadData[element.dataKey][i];
                            tableHtml += '<tr>';
                            element.rows[1].forEach(key => {
                                const actualKey = key.split('.').pop();
                                const cellValue = row[actualKey];
                                tableHtml += `<td>${cellValue !== undefined ? cellValue : ''}</td>`;
                            });
                            tableHtml += '</tr>';
                        }
                    }

                    tableHtml += `
            </tbody>
          </table>
        </div>`;
                    pageContent += tableHtml;
                    tableHeight += (endRow - startRow) * rowHeight; // Hitung tinggi tabel
                }

                yOffset = adjustedY + tableHeight + 20; // Update yOffset ke bawah tabel dengan jarak 20px
                return '';
            };

            const createNewPage = () => {
                htmlContent += `
    <div class="page">
      <div class="header">${headerContent}</div>
      ${pageContent}
      <div class="footer">${footerContent}</div>
    </div>`;
                pageContent = '';
                currentPage++;
                yOffset = 0; // Reset yOffset untuk halaman baru
            };

            this.elements.forEach(element => {
                if (element.y + element.height > currentPage * pageHeight) {
                    createNewPage();
                }
                addElementToPage(element);
                yOffset = Math.max(yOffset, element.y + element.height);
            });

            createNewPage(); // Tambahkan halaman terakhir

            htmlContent += `
  </body>
</html>`;

            this.sendToBackend(htmlContent);
        },
        async sendToBackend(htmlContent) {
            try {
                const response = await fetch('http://localhost:3000/generate-pdf', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'text/plain'  // Pastikan ini sesuai dengan yang backend harapkan
                    },
                    body: htmlContent,  // Kirim konten HTML sebagai teks biasa
                });

                if (response.ok) {
                    const blob = await response.blob();
                    const url = window.URL.createObjectURL(blob);
                    const a = document.createElement('a');
                    a.href = url;
                    a.download = 'generated.pdf';  // Nama default file PDF
                    document.body.appendChild(a);
                    a.click();
                    a.remove();  // Hapus elemen setelah diunduh
                } else {
                    console.error('Failed to generate PDF:', response.statusText);
                }
            } catch (error) {
                console.error('Error sending HTML to backend:', error);
            }
        },
        createHTML(htmlContent) {
            // Buat sebuah blob dari HTML content
            const blob = new Blob([htmlContent], { type: 'text/html' });

            // Buat URL objek dari blob
            const url = URL.createObjectURL(blob);

            // Buat elemen 'a' untuk mendownload file
            const a = document.createElement('a');
            a.href = url;
            a.download = 'report.html'; // Nama file HTML yang akan didownload
            document.body.appendChild(a); // Tambahkan elemen 'a' ke dalam body
            a.click(); // Trigger klik untuk memulai download
            document.body.removeChild(a); // Hapus elemen 'a' setelah download selesai

            // Bebaskan URL objek setelah tidak digunakan lagi
            URL.revokeObjectURL(url);
        },
        loadImage(src) {
            return new Promise((resolve, reject) => {
                const img = new Image();
                img.onload = () => resolve(img);
                img.onerror = reject;
                img.src = src;
            });
        },
        downloadPDF() {
            const doc = new jsPDF();
            this.elements.forEach((element) => {
                if (element.type === 'text') {
                    doc.text(element.x + this.elementMargin, element.y + this.elementMargin, element.content);
                }
                if (element.type === 'image') {
                    doc.addImage(element.src, 'JPEG', element.x + this.elementMargin, element.y + this.elementMargin, element.width, element.height);
                }
            });
            doc.save('report.pdf');
        },
        saveAsJSON() {
            // Clone elements untuk memodifikasi tanpa mengubah aslinya
            const layoutElements = this.elements.map(element => {
                const clonedElement = { ...element };
                if (clonedElement.type === 'table') {
                    // Hapus data payload jika ada
                    delete clonedElement.data;
                }
                return clonedElement;
            });

            // Simpan layout tanpa data payload
            const json = JSON.stringify(layoutElements, null, 2);
            const blob = new Blob([json], { type: "application/json" });
            saveAs(blob, "layout.json");
        },
        loadFromJSON(event) {
            const file = event.target.files[0];
            const reader = new FileReader();
            reader.onload = (e) => {
                const data = JSON.parse(e.target.result);

                // Verifikasi dan inisialisasi data untuk setiap elemen
                data.forEach(element => {
                    if (element.type === 'table') {
                        // Pastikan rows selalu memiliki dua baris: header dan placeholder
                        if (!element.rows) {
                            element.rows = [[], []];
                        } else {
                            if (!element.rows[0]) element.rows[0] = [];
                            if (!element.rows[1]) element.rows[1] = [];
                        }
                        // Pastikan data selalu diinisialisasi sebagai array
                        if (!element.data) {
                            element.data = [];
                        }
                    }
                });

                this.elements = data;
            };
            reader.readAsText(file);
        },
        extractKey(placeholder) {
            // Mengambil bagian terakhir dari placeholder sebagai key untuk data
            return placeholder.split('.').pop();
        },
        // Modifikasi render tabel untuk menggunakan extractKey
        renderTable(element) {
            return element.data.map((row) => {
                return element.rows[1].map((placeholder) => {
                    const key = this.extractKey(placeholder);
                    return row[key] !== undefined ? row[key] : 'Not Found';
                });
            });
        },
        addText() {
            const newText = {
                id: Date.now(),
                type: 'text',
                content: 'New Text',
                x: 20,
                y: 20,
                width: 150,
                height: 30,
                style: {},
                isTitle: false,
                isHeader: false,
                isBody: false,
                isFooter: false,
            };
            this.elements.push(newText);
        },
        addImage() {
            const newImage = {
                id: Date.now(),
                type: 'image',
                src: 'https://via.placeholder.com/150',
                x: 20,
                y: 20,
                width: 150,
                height: 150,
                style: {},
                isHeader: false,
                isBody: false,
                isFooter: false,
            };
            this.elements.push(newImage);
        },
        addTable() {
            const newTable = {
                id: Date.now(),
                type: 'table',
                rows: [
                    [], // Mulai dengan row header kosong
                    []  // Mulai dengan row placeholder kosong
                ],
                data: [], // Placeholder data untuk diisi dari payload JSON
                x: 20,
                y: 20,
                width: 600,
                height: 200,
                style: {},
                isHeader: false,
                isBody: false,
                isFooter: false,
            };
            this.elements.push(newTable);
        },
        addTableHeader() {
            if (this.selectedElement !== null && this.elements[this.selectedElement].type === 'table') {
                this.elements[this.selectedElement].rows[0].push("New Header");
            }
        },
        addTablePlaceholder() {
            if (this.selectedElement !== null && this.elements[this.selectedElement].type === 'table') {
                this.elements[this.selectedElement].rows[1].push("data.newField");
            }
        },
        selectElement(index) {
            this.selectedElement = index;
        },
        editText(index) {
            const element = this.elements[index];
            if (element.type === 'text') {
                this.$nextTick(() => {
                    const span = this.$el.querySelectorAll('.element span')[index];
                    if (span) {
                        span.focus();
                        span.addEventListener('blur', () => {
                            element.content = span.innerText;
                        });
                    } else {
                        console.error('Tidak dapat menemukan elemen span untuk fokus.');
                    }
                });
            }
        },
    },
};
</script>

<style scoped>
.header {
    position: relative;
    height: 100px;
    /* Tinggi header */
    border-bottom: 2px dashed #ccc;
    /* Garis batas bawah header */
    text-align: center;
    background-color: #f0f0f0;
}

.content-area {
    position: relative;
    flex-grow: 1;
    min-height: 400px;
    /* Minimal tinggi untuk area konten */
    padding: 10px;
    margin-bottom: 20px;
    /* Jarak dari footer */
}

.footer {
    position: relative;
    height: 100px;
    /* Tinggi footer */
    border-top: 2px dashed #ccc;
    /* Garis batas atas footer */
    text-align: center;
    background-color: #f0f0f0;
}

.report-builder {
    display: flex;
}

.menu {
    /* margin-right: 20px; */
}

.element {
    cursor: move;
}

.canvas-container {
    position: relative;
    /* width: 100%; */
    height: 100%;
    border: 1px solid #000;
    background: repeating-linear-gradient(45deg,
            #f0f0f0,
            #f0f0f0 10px,
            #ffffff 10px,
            #ffffff 20px);
    /* Border untuk menandai area canvas */
}

.no-grid {
    background: none !important;
    border: none !important;
}

.canvas {
    position: relative;
    width: 754px;
    /* Canvas A4 dengan margin 20px */
    height: 1083px;
    /* border: 1px solid #ccc; */
    background-color: #f9f9f9;
    background-image:
        linear-gradient(to right, #ccc 1px, transparent 1px),
        linear-gradient(to bottom, #ccc 1px, transparent 1px);
    background-size: 20px 20px;
    /* Ukuran grid sesuai gridSize */
    overflow: hidden;
}

.detail-form {
    margin-left: 20px;
    border: 1px solid #ccc;
    padding: 10px;
}

.element {
    cursor: move;
    box-sizing: border-box;
    border: 1px solid #000;
    /* Berikan border untuk setiap elemen */
}

.dragging {
    opacity: 0.5;
}

.table {
    width: 100%;
    border-collapse: collapse;
}

.table th,
.table td {
    border: 1px solid #000;
    /* Border for each cell */
    padding: 8px;
    text-align: center;
    text-transform: capitalize;
}

.table th {
    background-color: #f2f2f2;
}

.table-container {
    margin-top: 20px;
}

.delete-button {
    position: absolute;
    top: 5px;
    right: 5px;
    background-color: #ff4136;
    color: white;
    border: none;
    padding: 2px 5px;
    cursor: pointer;
    font-size: 12px;
}

.delete-button:hover {
    background-color: #d0342b;
}

@media print {
    .page {
        page-break-after: always;
    }

    table {
        page-break-inside: auto;
        border-collapse: collapse;
        width: 100%;
    }

    tr {
        page-break-inside: avoid;
        page-break-after: auto;
    }

    td,
    th {
        /* border: 1px solid #ddd; */
        padding: 8px;
    }
}
</style>