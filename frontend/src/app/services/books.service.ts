import { Injectable } from '@angular/core';
import { Book } from '../interfaces/book';

@Injectable({
  providedIn: 'root'
})
export class BooksService {

  constructor() { }

  dummyCovers = [
    'https://www.rockingbookcovers.com/wp-content/uploads/2021/01/No-One-Close.jpg',
    'https://www.rockingbookcovers.com/wp-content/uploads/2022/09/Dont-Let-Me-In.jpg',
    'https://www.rockingbookcovers.com/wp-content/uploads/2023/03/Dont-Go-There.jpg',
    'https://www.rockingbookcovers.com/wp-content/uploads/2024/01/When-I-Find-You.jpg',
    'https://www.rockingbookcovers.com/wp-content/uploads/2023/02/You-Shouldnt-Have.jpg',
  ]

  getDummyBooks(bookCount = 5): Book[] {
    return Array.from({ length: bookCount }, (_, index) => ({
      title: `Book ${index + 1}`,
      slag: `book-${index + 1}`,
      cover: this.dummyCovers[Math.floor(Math.random() * this.dummyCovers.length)]
    }))
  }

  getDummyBookBySlag(slag: string): Book | null {
    return {
      title: `${slag.replace(/-/g, ' ').toUpperCase()}`,
      slag: slag,
      cover: this.dummyCovers[Math.floor(Math.random() * this.dummyCovers.length)],
      short_description: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor.`,
      long_description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor. Cras elementum ultrices diam. Maecenas ligula massa, varius a, semper congue, euismod non, mi. Proin porttitor, orci nec nonummy molestie, enim est eleifend mi, non fermentum diam nisl sit amet erat. Duis semper. Duis arcu massa, scelerisque vitae, consequat in, pretium a, enim. Pellentesque congue. Ut in risus volutpat libero pharetra tempor. Cras elementum ultrices diam. Maecenas ligula massa, varius a, semper congue, euismod non, mi. Proin porttitor, orci nec nonummy molestie, enim est eleifend mi, non fermentum diam nisl sit amet erat. Duis semper. Duis arcu massa, scelerisque vitae, consequat in, pretium a, enim. Pellentesque congue. Ut in risus volutpat libero pharetra tempor.',
      author_avatar: 'https://www.picsum.photos/200',
      author_name: 'Author Name',
      author_short_bio: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor.',
      editors: ['Editor 1', 'Editor 2'],
      language: 'English',
      paperback: 'paper textured, matte finish, 300gsm, 5x8 inches, 300 pages',
    }
  }

  getBookBySlag(slag: string): Book | null {
    return this.getDummyBooks().find(book => book.slag === slag) || null;
  }
}
