const hello = document.getElementById('hello');
const skills = document.getElementById('skills');
const projects = document.getElementById('projects');
// const contacts = document.getElementById('contacts');

// const fonts = ['Ubuntu', 'Noto Sans', 'Serif', 'Verdana', 'FiraCode'];
const hello_list = ['Hello!', '你好！', 'नमस्ते!', '¡Hola!', 'Bonjour!', 'Привет!', 'Halo!','ハロー!','Hallo!', 'Kamusta!', 'toki a', 'Saluton!'];
const skills_list = ['My gizmos', '我的小玩意', 'मेरे उपकरण', 'mis artilugios', 'Mes gadgets', 'Мои вещицы', 'alat saya','私のギズモ','Meine Dinge', 'Aking mga gizmos', 'toki a', 'Saluton!', 'ilo mi', 'Miaj aparatoj'];
const projects_list = ['Projects', '项目', 'परियोजनाओं', 'Proyectos', 'Projets', 'Проекты', 'Proyek', 'プロジェクト', 'Projekte','Mga proyekto', 'pali', 'Projektoj'];
// const contacts_list = ['Contacts', 'Hola!', 'संपर्क', 'Contactos', 'Contacts', 'Контакты', 'Kontak', '連絡先','Kontakte', 'Mga contact', 'lipu mi', 'Kontaktoj'];




async function changing_hello(hello, hello_list, skills, skills_list, projects, projects_list, contacts, ccontacts_list) {
        while (true) {
                for (var i = 0; i < hello_list.length; i++) {
                        hello.innerHTML = hello_list[i];
                        skills.innerHTML = skills_list[i];
                        projects.innerHTML = projects_list[i];
                        // contacts.innerHTML = contacts_list[i];
                        // greeting.style.font_family = fonts[i];
                        await new Promise((r) => setTimeout(r, 1000));
                }
        }
}

changing_hello(hello, hello_list, skills, skills_list, projects, projects_list);
// changing_hello(hello, hello_list, skills, skills_list, projects, projects_list, contacts, contacts_list);
