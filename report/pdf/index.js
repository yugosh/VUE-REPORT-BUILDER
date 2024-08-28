const express = require('express');
const bodyParser = require('body-parser');
// const PdfPrinter = require('pdfmake');
// const vfsFonts = require('pdfmake/build/vfs_fonts');
const cors = require('cors');
const puppeteer = require('puppeteer');
const path = require('path');

const fs = require('fs');
// const path = require('path');

// const fonts = {
//     Roboto: {
//       normal: 'node_modules/pdfmake/build/vfs_fonts.js',
//     },
//   };

// const printer = new PdfPrinter(fonts);
const app = express();

app.use(cors());
app.use(bodyParser.json({ limit: '10mb' }));

app.post('/generate-pdf', async (req, res) => {
    const { htmlContent } = req.body;
  
    try {
      const browser = await puppeteer.launch();
      const page = await browser.newPage();
  
      // Load the HTML content into the page
      await page.setContent(htmlContent, { waitUntil: 'networkidle0' });
  
      // Generate PDF
      const pdfBuffer = await page.pdf({
        format: 'A4',
        printBackground: true,
      });
  
      await browser.close();
  
      // Step 4: Save the PDF to a local file to verify it's generated correctly
      const pdfPath = path.join(__dirname, 'test.pdf');
      fs.writeFileSync(pdfPath, pdfBuffer);
      console.log(`PDF saved to: ${pdfPath}`);
  
      // Step 5: Send the PDF as a response with proper headers
      res.setHeader('Content-Type', 'application/pdf');
      res.setHeader('Content-Disposition', 'attachment; filename=report.pdf');
      res.send(pdfBuffer);  // Mengirimkan PDF sebagai buffer
    } catch (error) {
      console.error('Error generating PDF:', error);
      res.status(500).send('Error generating PDF');
    }
  });
const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Server is running on port ${PORT}`);
});
