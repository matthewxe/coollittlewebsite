const hello = document.getElementById('hello');
// const skills = document.getElementById('skills');
// const projects = document.getElementById('projects');
// const contacts = document.getElementById('contacts');

// const fonts = ['Ubuntu', 'Noto Sans', 'Serif', 'Verdana', 'FiraCode'];
const hello_list = ['Hello!', '你好！', 'नमस्ते!', '¡Hola!', 'Bonjour!', 'Привет!', 'Halo!','ハロー!','Hallo!', 'Kamusta!', 'toki a', 'Saluton!'];
// const skills_list = ['Hello!', 'Hola!', 'Hallo', 'Kamusta!', 'こんにちわ', 'toki a', '你好'];
// const projects_list = ['Hello!', 'Hola!', 'Hallo', 'Kamusta!', 'こんにちわ', 'toki a', '你好'];
// const contacts_list = ['Hello!', 'Hola!', 'Hallo', 'Kamusta!', 'こんにちわ', 'toki a', '你好'];




async function changing_hello(hello, hello_list, skills, skills_list, projects, projects_list) {
        while (true) {
                for (var i = 0; i < hello_list.length; i++) {
                        hello.innerHTML = hello_list[i];
                        // skills.innerHTML = skills_list[i];
                        // greeting.style.font_family = fonts[i];
                        await new Promise((r) => setTimeout(r, 1000));
                }
        }
}

changing_hello(hello, hello_list);
// changing_hello(hello, hello_list, skills, skills_list, projects, projects_list, contacts, contacts_list);
