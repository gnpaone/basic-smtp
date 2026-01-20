document.addEventListener('DOMContentLoaded', () => {
    const emailList = document.getElementById('email-list');
    const emailDetail = document.getElementById('email-detail');
    const noSelection = document.getElementById('no-selection');
    const toggleHeadersBtn = document.getElementById('toggle-headers-btn');
    const detailHeaders = document.getElementById('detail-headers');

    let currentEmails = [];
    let selectedEmailId = null;

    const detailSubject = document.getElementById('detail-subject');
    const detailFrom = document.getElementById('detail-from');
    const detailTo = document.getElementById('detail-to');
    const detailTime = document.getElementById('detail-time');
    const detailBody = document.getElementById('detail-body');
    const detailHeadersPre = document.getElementById('detail-headers');

    function fetchEmails() {
        console.log("Fetching emails...");
        fetch('/api/emails')
            .then(response => response.json())
            .then(data => {
                if (!Array.isArray(data)) {
                    data = [];
                }
                renderEmailList(data);
                currentEmails = data;
            })
            .catch(error => console.error('Error fetching emails:', error));
    }

    function renderEmailList(emails) {
        if (emails.length === 0) {
            emailList.innerHTML = '<div class="empty-state">Waiting for emails...</div>';
            return;
        }

        emailList.innerHTML = '';
        emails.forEach(email => {
            const item = document.createElement('div');
            item.className = 'email-item';
            if (email.id === selectedEmailId) {
                item.classList.add('selected');
            }

            const date = new Date(email.created_at);
            const timeStr = date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });

            item.innerHTML = `
                <div class="email-item-header">
                    <span class="email-sender">${escapeHtml(email.from)}</span>
                    <span class="email-time">${timeStr}</span>
                </div>
                <h4 class="email-subject">${escapeHtml(email.subject || '(No Subject)')}</h4>
            `;

            item.addEventListener('click', () => selectEmail(email));
            emailList.appendChild(item);
        });
    }

    function selectEmail(email) {
        selectedEmailId = email.id;

        document.querySelectorAll('.email-item').forEach(el => el.classList.remove('selected'));
        renderEmailList(currentEmails); 

        noSelection.classList.add('hidden');
        emailDetail.classList.remove('hidden');

        detailSubject.textContent = email.subject || '(No Subject)';
        detailFrom.textContent = `From: ${email.from}`;
        detailTo.textContent = email.to ? email.to.join(', ') : '';
        detailTime.textContent = new Date(email.created_at).toLocaleString();
        detailBody.textContent = email.body;
        detailHeadersPre.textContent = email.headers;
    }

    toggleHeadersBtn.addEventListener('click', () => {
        detailHeaders.classList.toggle('hidden');
    });

    function escapeHtml(text) {
        if (!text) return '';
        const map = {
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            '"': '&quot;',
            "'": '&#039;'
        };
        return text.replace(/[&<>"']/g, function(m) { return map[m]; });
    }

    fetchEmails();
    setInterval(fetchEmails, 3000);
});
