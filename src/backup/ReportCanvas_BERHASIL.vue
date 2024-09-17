<template>
    <div class="report-builder">
        <!-- Sidebar Menu -->
        <div class="menu">
            <button @click="addText">Add Text</button>
            <button @click="addImage">Add Image</button>
            <button @click="addTable">Add Table</button>
            <input type="file" @change="loadDataPayload" accept=".json" />
        </div>

        <!-- Canvas Area -->
        <div class="canvas-container">
            <div class="header">
                <h2>Header Content</h2>
            </div>
            <div class="canvas" @dragover.prevent @drop="onDrop" ref="canvas">
                <draggable v-model="elements" :itemKey="getItemKey"
                    :options="{ group: 'elements', ghostClass: 'dragging', delay: 0, delayOnTouchOnly: true }">
                    <template v-slot:item="{ element, index }">
                        <div draggable="true" :key="getItemKey(element, index)" @dragstart="onDragStart($event, index)"
                            @click="selectElement(index)" :style="getElementStyle(element)" class="element">
                            <span v-if="element.type === 'text'" @dblclick="editText(index)" contenteditable="true">{{
                element.content }}</span>
                            <img v-if="element.type === 'image'" :src="element.src" alt="Image" />
                            <!-- <table v-if="element.type === 'table'" :style="getTableStyle(element)">
                                <thead>
                                    <tr>
                                        <th v-for="(cell, cellIndex) in element.rows[0]" :key="cellIndex">{{ cell }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody>
                                    Jika tidak ada data, tampilkan placeholder
                                    <tr v-if="element.data.length === 0">
                                        <td v-for="(placeholder, cellIndex) in element.rows[1]" :key="cellIndex">
                                            {{ placeholder }}
                                        </td>
                                    </tr>
                                    Render data dari payload JSON
                                    <tr v-for="(row, rowIndex) in element.data" :key="rowIndex">
                                        <td v-for="(placeholder, cellIndex) in element.rows[1]" :key="cellIndex">
                                            {{ row[extractKey(placeholder)] !== undefined ? row[extractKey(placeholder)]
                : 'Not Found' }}
                                        </td>
                                    </tr>
                                </tbody>
                            </table> -->

                            <table v-if="element.type === 'table'" :style="getTableStyle(element)">
                                <thead>
                                    <tr>
                                        <th v-for="(cell, cellIndex) in element.rows[0]" :key="cellIndex">{{ cell }}
                                        </th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr v-for="(row, rowIndex) in getTableData(element.dataKey)" :key="rowIndex">
                                        <td v-for="(placeholder, cellIndex) in element.rows[1]" :key="cellIndex">
                                            {{ row[extractKey(placeholder)] || 'N/A' }}
                                        </td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    </template>
                </draggable>
            </div>
            <div class="footer">
                <h3>Footer Content</h3>
            </div>
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
                </div>
            </div>
        </div>

        <button @click="saveLayout">Save Layout</button>
        <button @click="generatePDF">Download PDF</button>
        <button @click="saveAsJSON">Save as JSON</button>
        <input type="file" @change="loadFromJSON" accept=".json" />
    </div>
</template>

<script>
import draggable from 'vuedraggable';
import jsPDF from 'jspdf';
import { saveAs } from 'file-saver';
import html2canvas from 'html2canvas';

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
            payloadData: [],
        };
    },
    watch: {
        payloadData: {
            handler(newData) {
                console.log('Payload data updated:', newData);
                this.updateAllTables();
            },
            deep: true
        }
    },
    methods: {
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
                border: '1px solid #000', // Tambahkan border ke elemen
                boxSizing: 'border-box',
            };
        },
        getTableStyle(element) {
            return {
                width: `${element.width}px`,
                height: `${element.height}px`,
                border: '1px solid #000',
                borderCollapse: 'collapse',
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
            const canvasContainer = this.$refs.canvas;

            //   // Tambahkan kelas khusus untuk menghapus grid atau background
            canvasContainer.classList.add('no-grid');

            html2canvas(canvasContainer, {
                useCORS: true,
                allowTaint: false
            }).then((canvas) => {
                canvas.toBlob((blob) => {
                    const reader = new FileReader();
                    reader.onload = function (event) {
                        const imgData = event.target.result;
                        const pdf = new jsPDF({
                            orientation: 'portrait',
                            unit: 'px',
                            format: [canvas.width, canvas.height]
                        });

                        try {
                            pdf.addImage(imgData, 'PNG', 0, 0, canvas.width, canvas.height);
                            canvasContainer.classList.remove('no-grid');

                            pdf.save('canvas-report.pdf');
                        } catch (error) {
                            console.error("Error generating PDF: ", error);
                        }
                    };
                    reader.readAsDataURL(blob);
                });
            }).catch(error => {
                console.error("Error generating PDF: ", error);
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
        // loadDataPayload(event) {
        //     const file = event.target.files[0];
        //     const reader = new FileReader();
        //     reader.onload = (e) => {
        //         const data = JSON.parse(e.target.result);
        //         this.elements.forEach((element) => {
        //             if (element.type === 'table') {
        //                 element.data = data;
        //             }
        //         });
        //     };
        //     reader.readAsText(file);
        // },
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
                width: 100,
                height: 30,
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
    margin-right: 20px;
}

.element {
    cursor: move;
}

.canvas-container {
    position: relative;
    width: 100%;
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
    border: 1px solid #ccc;
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

table {
    width: 100%;
    height: 100%;
    border-collapse: collapse;
}

th,
td {
    border: 1px solid #000;
    text-align: center;
    padding: 4px;
}
</style>