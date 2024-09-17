<template>
    <div class="report-builder">
        <!-- Sidebar Menu -->
        <div class="menu">
            <button @click="addText">Add Text</button>
            <button @click="addImage">Add Image</button>
            <button @click="addTable">Add Table</button>

            <button @click="addField">Field</button>

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
                            <span v-if="element.type === 'text'" @dblclick="editText(index)" contenteditable="true">
                                <!-- {{ element.content }} -->
                                {{ parseTextContent(element.content) }}
                            </span>
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

                            <button @click.stop="deleteElement(index)" class="delete-button">Delete</button>
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
    methods: {
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
        // generatePDF() {
        //     if (!this.payloadData) {
        //         console.error('Payload data is not available');
        //         return;
        //     }

        //     const pdf = new jsPDF('p', 'mm', 'a4');
        //     const pageHeight = pdf.internal.pageSize.getHeight();
        //     let startY = 30; // Awal Y untuk data tabel

        //     const addHeader = () => {
        //         pdf.setFontSize(14);
        //         pdf.text("Header Content", pdf.internal.pageSize.getWidth() / 2, 10, { align: 'center' });
        //     };

        //     const addFooter = (pageNumber) => {
        //         pdf.setFontSize(10);
        //         pdf.text(`Page ${pageNumber}`, pdf.internal.pageSize.getWidth() / 2, pageHeight - 10, { align: 'center' });
        //     };

        //     let pageNumber = 1;
        //     addHeader();
        //     addFooter(pageNumber);

        //     this.elements.forEach(element => {
        //         if (element.type === 'table') {
        //             const tableData = element.dataKey ? this.payloadData[element.dataKey] : [];

        //             // Format data untuk digunakan dengan jsPDF autotable
        //             const rows = tableData.map(row => [row.nama, row.gaji.toString()]);
        //             const headers = [["Nama", "Gaji"]];

        //             pdf.autoTable({
        //                 head: headers,
        //                 body: rows,
        //                 startY: startY,
        //                 theme: 'grid',
        //                 styles: {
        //                     fontSize: 10,
        //                     cellPadding: 3,
        //                     valign: 'middle',
        //                     halign: 'center', // Horizontal alignment for all columns
        //                     lineColor: [44, 62, 80],
        //                     lineWidth: 0.75,
        //                 },
        //                 headStyles: {
        //                     fillColor: [52, 152, 219],
        //                     textColor: [255, 255, 255],
        //                     fontStyle: 'bold',
        //                 },
        //                 footStyles: {
        //                     fillColor: [52, 152, 219],
        //                     textColor: [255, 255, 255],
        //                     fontStyle: 'bold',
        //                 },
        //                 alternateRowStyles: {
        //                     fillColor: [240, 240, 240],
        //                 },
        //                 didDrawPage: () => {
        //                     // Add footer with page number
        //                     addFooter(pageNumber);
        //                     pageNumber++;
        //                 },
        //                 margin: { top: 40 }, // Adjust top margin to prevent overlap with header
        //             });

        //             startY = pdf.lastAutoTable.finalY + 10; // Update startY to be after the table
        //         }
        //     });

        //     // Simpan PDF
        //     pdf.save('canvas-report.pdf');
        // },
        // generatePDF() {
        //     if (!this.payloadData) {
        //         console.error('Payload data is not available');
        //         return;
        //     }

        //     const pdf = new jsPDF('p', 'mm', 'a4');
        //     const pageHeight = pdf.internal.pageSize.getHeight();
        //     let startY = 30; // Awal Y untuk data tabel

        //     const addHeader = () => {
        //         pdf.setFontSize(14);
        //         pdf.text("Header Content", pdf.internal.pageSize.getWidth() / 2, 10, { align: 'center' });
        //     };

        //     const addFooter = (pageNumber) => {
        //         pdf.setFontSize(10);
        //         pdf.text(`Page ${pageNumber}`, pdf.internal.pageSize.getWidth() / 2, pageHeight - 10, { align: 'center' });
        //     };

        //     let pageNumber = 1;
        //     addHeader();
        //     addFooter(pageNumber);

        //     const headerMapping = {
        //         "Tanggal": "fromDate",
        //         "Jam Masuk": "checkin",
        //         "Jam Pulang": "checkout",
        //         "Keluar Istirahat": "breakedout",
        //         "Kembali Istirahat": "breakedin",
        //         "Keterangan Sistem": "holidayexplanation",
        //         "Keterangan": "remark"
        //     };

        //     this.elements.forEach(element => {
        //         switch (element.type) {
        //             case 'table': {
        //                 const tableData = element.dataKey ? this.payloadData[element.dataKey] : [];

        //                 const headers = Object.keys(headerMapping);

        //                 const rows = tableData.map(row => {
        //                     return headers.map(header => {
        //                         const jsonKey = headerMapping[header];
        //                         return row[jsonKey] ? row[jsonKey].toString() : 'N/A';
        //                     });
        //                 });

        //                 pdf.autoTable({
        //                     head: [headers],
        //                     body: rows,
        //                     startY: startY,
        //                     theme: 'grid',
        //                     styles: {
        //                         fontSize: 10,
        //                         cellPadding: 3,
        //                         valign: 'middle',
        //                         halign: 'center',
        //                         lineColor: [44, 62, 80],
        //                         lineWidth: 0.75,
        //                     },
        //                     headStyles: {
        //                         fillColor: [52, 152, 219],
        //                         textColor: [255, 255, 255],
        //                         fontStyle: 'bold',
        //                     },
        //                     footStyles: {
        //                         fillColor: [52, 152, 219],
        //                         textColor: [255, 255, 255],
        //                         fontStyle: 'bold',
        //                     },
        //                     alternateRowStyles: {
        //                         fillColor: [240, 240, 240],
        //                     },
        //                     didDrawPage: () => {
        //                         addFooter(pageNumber);
        //                         pageNumber++;
        //                     },
        //                     margin: { top: 40 },
        //                 });

        //                 startY = pdf.lastAutoTable.finalY + 10;
        //                 break;
        //             }

        //             case 'text': {
        //             if (element.content) {
        //                 pdf.setFontSize(element.fontSize || 12);
        //                 pdf.text(element.content, element.x || 10, startY);
        //                 startY += 10; // Adjust the Y position for the next element
        //             }
        //             break;
        //         }

        //         case 'field': {
        //             if (element.field && this.payloadData[element.field]) {
        //                 pdf.setFontSize(element.fontSize || 12);
        //                 pdf.text(this.payloadData[element.field].toString(), element.x || 10, startY);
        //                 startY += 10; // Adjust the Y position for the next element
        //             }
        //             break;
        //         }

        //         case 'image': {
        //             if (element.src) {
        //                 const img = new Image();
        //                 img.src = element.src;
        //                 img.onload = () => {
        //                     const imgWidth = element.width || 150;
        //                     const imgHeight = element.height || 150;
        //                     const scaleWidth = imgWidth / img.width;
        //                     const scaleHeight = imgHeight / img.height;
        //                     const scale = Math.min(scaleWidth, scaleHeight); // Scaling to maintain aspect ratio

        //                     pdf.addImage(
        //                         img,
        //                         'JPEG',
        //                         element.x || 10,
        //                         startY,
        //                         img.width * scale,
        //                         img.height * scale
        //                     );
        //                     startY += img.height * scale + 10; // Adjust the Y position for the next element
        //                 };
        //             }
        //             break;
        //         }

        //             default:
        //                 console.warn('Unknown element type:', element.type);
        //                 break;
        //         }
        //     });

        //     pdf.save('canvas-report.pdf');
        // },


        //GENERATEPDF DIBAWAH HAMPIR MIRIP
        // async generatePDF() {
        //     if (!this.payloadData) {
        //         console.error('Payload data is not available');
        //         return;
        //     }

        //     const pdf = new jsPDF('p', 'pt', 'a4');
        //     const pageWidth = pdf.internal.pageSize.getWidth();
        //     const pageHeight = pdf.internal.pageSize.getHeight();
        //     const scaleFactor = pageWidth / this.canvasWidth;

        //     const addHeader = () => {
        //         pdf.setFontSize(14);
        //         pdf.text("Header Content", pageWidth / 2, 20, { align: 'center' });
        //     };

        //     const addFooter = (pageNumber) => {
        //         pdf.setFontSize(10);
        //         pdf.text(`Page ${pageNumber}`, pageWidth / 2, pageHeight - 10, { align: 'center' });
        //     };

        //     let pageNumber = 1;
        //     addHeader();
        //     addFooter(pageNumber);

        //     const headerMapping = {
        //         "Tanggal": "fromDate",
        //         "Jam Masuk": "checkin",
        //         "Jam Pulang": "checkout",
        //         "Keluar Istirahat": "breakedout",
        //         "Kembali Istirahat": "breakedin",
        //         "Keterangan Sistem": "holidayexplanation",
        //         "Keterangan": "remark"
        //     };

        //     const translateY = (y) => y * scaleFactor;

        //     // Sort elements by their y position
        //     this.elements.sort((a, b) => a.y - b.y);

        //     for (const element of this.elements) {
        //         const x = element.x * scaleFactor;
        //         let y = translateY(element.y);

        //         // Check if we need to add a new page
        //         if (y > pageHeight - 40) {
        //             pdf.addPage();
        //             pageNumber++;
        //             addHeader();
        //             addFooter(pageNumber);
        //             y = 40; // Reset y for the new page
        //         }

        //         switch (element.type) {
        //             case 'text':
        //             case 'field': {
        //                 pdf.setFontSize(element.fontSize || 12);
        //                 const text = element.type === 'text' ? element.content :
        //                     (this.payloadData[element.field] || '').toString();
        //                 pdf.text(text, x, y);
        //                 break;
        //             }
        //             case 'image': {
        //                 if (element.src) {
        //                     try {
        //                         const img = await this.loadImage(element.src);
        //                         const imgWidth = element.width * scaleFactor;
        //                         const imgHeight = element.height * scaleFactor;
        //                         pdf.addImage(img, 'PNG', x, y, imgWidth, imgHeight);
        //                     } catch (error) {
        //                         console.error('Error loading image:', error);
        //                     }
        //                 }
        //                 break;
        //             }
        //             case 'table': {
        //                 // Implement table rendering here
        //                 const tableData = element.dataKey ? this.payloadData[element.dataKey] : [];
        //                 const headers = Object.keys(headerMapping);
        //                 const rows = tableData.map(row =>
        //                     headers.map(header => {
        //                         const jsonKey = headerMapping[header];
        //                         return row[jsonKey] ? row[jsonKey].toString() : 'N/A';
        //                     })
        //                 );

        //                 pdf.autoTable({
        //                     head: [headers],
        //                     body: rows,
        //                     startY: y,
        //                     theme: 'grid',
        //                     styles: {
        //                         fontSize: 10,
        //                         cellPadding: 3,
        //                         valign: 'middle',
        //                         halign: 'center',
        //                         lineColor: [44, 62, 80],
        //                         lineWidth: 0.75,
        //                     },
        //                     headStyles: {
        //                         fillColor: [52, 152, 219],
        //                         textColor: [255, 255, 255],
        //                         fontStyle: 'bold',
        //                     },
        //                     footStyles: {
        //                         fillColor: [52, 152, 219],
        //                         textColor: [255, 255, 255],
        //                         fontStyle: 'bold',
        //                     },
        //                     alternateRowStyles: {
        //                         fillColor: [240, 240, 240],
        //                     },
        //                     didDrawPage: () => {
        //                         addFooter(pageNumber);
        //                         pageNumber++;
        //                         addHeader();
        //                     },
        //                     margin: { top: 40 },
        //                 });

        //                 // Update y to avoid overlap with the following content
        //                 y = pdf.lastAutoTable.finalY + 10;

        //                 break;
        //             }
        //             default: {
        //                 console.warn('Unknown element type:', element.type);
        //                 break;
        //             }
        //         }
        //     }

        //     pdf.save('canvas-report.pdf');
        // },
        async generatePDF2() {
            if (!this.payloadData) {
                console.error('Payload data is not available');
                return;
            }

            const pdf = new jsPDF('p', 'pt', 'a4');
            const pageWidth = pdf.internal.pageSize.getWidth();
            const pageHeight = pdf.internal.pageSize.getHeight();
            const scaleFactor = pageWidth / this.canvasWidth;
            let currentY = 40; // Start Y position after header
            let pageNumber = 1;

            const addHeader = () => {
                pdf.setFontSize(14);
                pdf.text("Header Content", pageWidth / 2, 20, { align: 'center' });
                currentY = 40; // Reset Y after header
            };

            const addFooter = () => {
                pdf.setFontSize(10);
                pdf.text(`Page ${pageNumber}`, pageWidth / 2, pageHeight - 10, { align: 'center' });
                pageNumber++;
            };

            const checkPageOverflow = (additionalHeight) => {
                if (currentY + additionalHeight > pageHeight - 40) {
                    pdf.addPage();
                    addHeader();
                    addFooter();
                    currentY = 40;
                }
            };

            const addText = (content, x, y, fontSize = 12) => {
                checkPageOverflow(20);
                pdf.setFontSize(fontSize);
                pdf.text(content, x * scaleFactor, currentY);
                currentY += 20; // Adjust according to font size
            };

            const addImage = async (src, x, y, width, height) => {
                const imgWidth = width * scaleFactor;
                const imgHeight = height * scaleFactor;
                checkPageOverflow(imgHeight);
                try {
                    const img = await this.loadImage(src);
                    pdf.addImage(img, 'PNG', x * scaleFactor, currentY, imgWidth, imgHeight);
                    currentY += imgHeight + 10;
                } catch (error) {
                    console.error('Error loading image:', error);
                }
            };

            const addTable = (data, x, y, headers) => {
                const tableData = data.map(row =>
                    headers.map(header => row[header] ? row[header].toString() : 'N/A')
                );
                checkPageOverflow(20 * tableData.length); // Approximation for table height
                pdf.autoTable({
                    head: [headers],
                    body: tableData,
                    startY: currentY,
                    margin: { left: x * scaleFactor },
                    theme: 'grid',
                    styles: {
                        fontSize: 10,
                        cellPadding: 3,
                        valign: 'middle',
                        halign: 'center',
                        lineColor: [44, 62, 80],
                        lineWidth: 0.75,
                    },
                    headStyles: {
                        fillColor: [52, 152, 219],
                        textColor: [255, 255, 255],
                        fontStyle: 'bold',
                    },
                    footStyles: {
                        fillColor: [52, 152, 219],
                        textColor: [255, 255, 255],
                        fontStyle: 'bold',
                    },
                    alternateRowStyles: {
                        fillColor: [240, 240, 240],
                    },
                    didDrawPage: () => {
                        addFooter();
                        addHeader();
                    }
                });
                currentY = pdf.lastAutoTable.finalY + 10;
            };

            addHeader();
            addFooter();

            const headerMapping = {
                "Tanggal": "fromDate",
                "Jam Masuk": "checkin",
                "Jam Pulang": "checkout",
                "Keluar Istirahat": "breakedout",
                "Kembali Istirahat": "breakedin",
                "Keterangan Sistem": "holidayexplanation",
                "Keterangan": "remark"
            };

            for (const element of this.elements) {
                const x = element.x * scaleFactor;
                let y = element.y * scaleFactor;

                switch (element.type) {
                    case 'text':
                    case 'field':{
                        const textContent = element.type === 'text' ? element.content :
                            (this.payloadData[element.field] || '').toString();
                        addText(textContent, x, y, element.fontSize || 12);
                        break;}
                    case 'image':{
                        await addImage(element.src, x, y, element.width, element.height);
                        break;}
                    case 'table':{
                        const tableData = element.dataKey ? this.payloadData[element.dataKey] : [];
                        const headers = Object.keys(headerMapping);
                        addTable(tableData, x, y, headers);
                        break;}
                    default:{
                        console.warn('Unknown element type:', element.type);
                        break;}
                }
            }

            pdf.save('canvas-report.pdf');
        },
        async generatePDF() {
    if (!this.payloadData) {
        console.error('Payload data is not available');
        return;
    }

    const pdf = new jsPDF('p', 'pt', 'a4');
    const pageWidth = pdf.internal.pageSize.getWidth();
    const pageHeight = pdf.internal.pageSize.getHeight();
    let pageNumber = 1;

    const addHeader = () => {
        pdf.setFontSize(14);
        pdf.text("Header Content", pageWidth / 2, 20, { align: 'center' });
    };

    const addFooter = () => {
        pdf.setFontSize(10);
        pdf.text(`Page ${pageNumber}`, pageWidth / 2, pageHeight - 10, { align: 'center' });
        pageNumber++;
    };

    addHeader();
    addFooter();

    // Posisi elemen-elemen, diatur secara manual
    const elements = [
        { type: 'text', content: 'Laporan Abnormal', x: 50, y: 50, fontSize: 14 },
        { type: 'text', content: '28382183374534', x: 50, y: 70, fontSize: 12 },
        { type: 'text', content: 'From Date : 09/09/2019 to : 10/10/2020', x: 50, y: 90, fontSize: 12 },
        { type: 'image', src: 'path/to/image.png', x: 50, y: 110, width: 100, height: 100 },
        { type: 'text', content: 'NIK', x: 50, y: 220, fontSize: 12 },
        { type: 'text', content: 'NAMA', x: 50, y: 240, fontSize: 12 },
        {
            type: 'table', 
            dataKey: 'laporanabnormal', 
            x: 50, 
            y: 260,
            headers: ["Tanggal", "Jam Masuk", "Jam Pulang", "Keluar Istirahat", "Kembali Istirahat", "Keterangan Sistem", "Keterangan"]
        }
    ];

    for (const element of elements) {
        const x = element.x;
        let y = element.y;

        // Check if we need to add a new page
        if (y > pageHeight - 60) {
            pdf.addPage();
            addHeader();
            addFooter();
            y = 50; // Reset y for the new page
        }

        switch (element.type) {
            case 'text': {
                pdf.setFontSize(element.fontSize || 12);
                pdf.text(element.content, x, y);
                break;
            }
            case 'image': {
                try {
                    const img = await this.loadImage(element.src);
                    const imgWidth = element.width;
                    const imgHeight = element.height;
                    pdf.addImage(img, 'PNG', x, y, imgWidth, imgHeight);
                } catch (error) {
                    console.error('Error loading image:', error);
                }
                break;
            }
            case 'table': {
                const tableData = this.payloadData[element.dataKey] || [];
                const rows = tableData.map(row =>
                    element.headers.map(header => row[header.toLowerCase()] || 'N/A')
                );

                pdf.autoTable({
                    head: [element.headers],
                    body: rows,
                    startY: y,
                    theme: 'grid',
                    styles: {
                        fontSize: 10,
                        cellPadding: 3,
                        valign: 'middle',
                        halign: 'center',
                        lineColor: [44, 62, 80],
                        lineWidth: 0.75,
                    },
                    headStyles: {
                        fillColor: [52, 152, 219],
                        textColor: [255, 255, 255],
                        fontStyle: 'bold',
                    },
                    footStyles: {
                        fillColor: [52, 152, 219],
                        textColor: [255, 255, 255],
                        fontStyle: 'bold',
                    },
                    alternateRowStyles: {
                        fillColor: [240, 240, 240],
                    },
                    margin: { left: x },
                });
                break;
            }
            default: {
                console.warn('Unknown element type:', element.type);
                break;
            }
        }
    }

    pdf.save('canvas-report-abnormal.pdf');
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