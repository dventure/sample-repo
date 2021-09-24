import imaplib
import email
import os
from reportlab.pdfgen import canvas


host = 'imap.gmail.com'
user = 'test1day2020@gmail.com'
pwd = 'strongword'

# Getting unread mail from inbox
def get_inbox():
    mail = imaplib.IMAP4_SSL(host)
    mail.login(user,pwd)
    mail.select('inbox')
    _, search_data = mail.search(None, 'UNSEEN')
    my_message = []

    for num in search_data[0].split():
        email_data = {}
        _, data = mail.fetch(num, 'RFC822')
        _,b = data[0]
        msg = email.message_from_bytes(b)
        
        for  header in ['subject', 'to', 'from', 'date']:
            email_data[header] = msg[header]
        file_name=' '    
        for part in msg.walk():            
            if part.get_content_type() == 'text/plain':
                body = part.get_payload(decode=True)
                email_data['body'] = body.decode()
##            elif part.get_content_type() == 'text/html':
##                html_body = part.get_payload(decode=True)
##                email_data['html_body'] = html_body.decode()
            elif part.get_filename() != None:
                file_name = part.get_filename()
                fp = open(file_name, 'wb')
                fp.write(part.get_payload(decode=True))
                fp.close()
                
            #print(file_name)
            email_data['attach'] = file_name   
        my_message.append(email_data)
   # print(my_message)
    return my_message        

#Generating pdf    
def gen_pdf(msg):
    sub_title = msg['subject']
    name = msg['date']+'.pdf'
    frm = msg['from'].split('<')[0]
    title = 'Email - ' + frm.title()
    body = msg['body']
    attachment = msg['attach']
    #textLines = ['Hello how', 'are', 'you']
    textLines = body.splitlines()
    #img = 'img.jpeg'
    img = 'https://www.howtogeek.com/wp-content/uploads/2019/03/gmail-1.png?width=250&trim=1,1&bg-color=000&pad=1,1'
    
    pdf = canvas.Canvas(name)
    
    pdf.setTitle(title)
    pdf.setFont('Courier', 30)
    pdf.drawCentredString(300,750, title)
    
    pdf.setFont('Courier', 25)
    pdf.setFillColorRGB(0,0,256)
    pdf.drawCentredString(290,700,sub_title)
    
    pdf.line(150, 650, 550, 650)
    
##    for font in pdf.getAvailableFonts():
##        print(font)
    from reportlab.lib import colors
    text = pdf.beginText(100, 600)
    text.setFont('Courier', 20)
    text.setFillColor(colors.red)
    
    for line in textLines:
        text.textLine(line)
        #print (line)
    pdf.drawText(text)

    #Attachment
    #text = pdf.beginText(100, 200)
    text.setFont('Courier', 20)
    text.setFillColor(colors.blue)
    pdf.drawCentredString(290,400,attachment)
    
    pdf.drawInlineImage(img, 150, 50)
    pdf.save()
    print("Generated pdf", name)

# Main flow
my_msgs = get_inbox()
##my_msgs = [{'subject': 'This is a subject', 'to': 'test1day2020@gmail.com', 'from': 'test account <test1day2020@gmail.com>', 'date': 'Wed, 15 Sep 2021 00:18:40 +0530', 'body': 'Hello,\r\n\r\nThis. is my body!\r\n\r\nThanks,\r\nGijoy\r\n', 'html_body': '<div dir="ltr">Hello,<div><br></div><div>This. is my body!</div><div><br></div><div>Thanks,</div><div>Gijoy</div></div>\r\n'},
##{'subject': 'New Subject', 'to': 'test1day2020@gmail.com', 'from': 'test account <test1day2020@gmail.com>', 'date': 'Wed, 15 Sep 2021 19:57:39 +0530', 'body': 'THis is the new email\r\nYou can check this mail\r\n\r\nThanks\r\nTest\r\n', 'html_body': '<div dir="ltr">THis is the new email<div>You can check this mail</div><div><br></div><div>Thanks</div><div>Test</div></div>\r\n'}]
filter_param = 'test1day2020'
for each_msg in my_msgs:
        if filter_param in each_msg['from']:
            gen_pdf(each_msg)

