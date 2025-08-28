# Advanced Tiptap Editor Installation and Usage Guide

## 1. Package Installation

```bash
cd web
npm install @tiptap/vue-3 @tiptap/pm @tiptap/starter-kit @tiptap/extension-placeholder @tiptap/extension-link @tiptap/extension-image @tiptap/extension-text-align @tiptap/extension-underline @tiptap/extension-text-style @tiptap/extension-color @tiptap/extension-highlight @tiptap/extension-task-list @tiptap/extension-task-item @tiptap/extension-table @tiptap/extension-table-row @tiptap/extension-table-cell @tiptap/extension-table-header
```

## 2. Package List to Install

### Core Packages

- `@tiptap/vue-3`: Tiptap component for Vue 3
- `@tiptap/pm`: ProseMirror (Tiptap foundation)
- `@tiptap/starter-kit`: Basic features

### 확장 기능

- `@tiptap/extension-placeholder`: 플레이스홀더 텍스트
- `@tiptap/extension-link`: 링크 기능
- `@tiptap/extension-image`: 이미지 삽입
- `@tiptap/extension-text-align`: 텍스트 정렬
- `@tiptap/extension-underline`: 밑줄
- `@tiptap/extension-text-style`: 텍스트 스타일
- `@tiptap/extension-color`: 텍스트 색상
- `@tiptap/extension-highlight`: 하이라이트
- `@tiptap/extension-task-list`: 체크리스트
- `@tiptap/extension-task-item`: 체크리스트 항목
- `@tiptap/extension-table`: 테이블
- `@tiptap/extension-table-row`: 테이블 행
- `@tiptap/extension-table-cell`: 테이블 셀
- `@tiptap/extension-table-header`: 테이블 헤더
- `@tiptap/extension-code-block`: 코드 블록

## 3. 고급 기능

### 🎨 텍스트 서식

- **굵게** (Ctrl+B)
- _기울임_ (Ctrl+I)
- <u>밑줄</u>
- ~~취소선~~
- `인라인 코드`

### 🎨 색상 및 하이라이트

- 텍스트 색상 변경
- 배경색 하이라이트
- 다양한 색상 팔레트

### 📋 목록 및 체크리스트

- 글머리 기호 목록
- 번호 매기기 목록
- 체크리스트 (체크박스)

### 📊 테이블

- 테이블 삽입
- 행/열 조정
- 테이블 헤더
- 셀 선택 및 편집

### 💻 코드 블록

- 기본 코드 블록 지원
- 다크 테마 코드 블록

### 🔗 링크 및 미디어

- 링크 삽입
- 이미지 삽입
- URL 입력 다이얼로그

### 📐 정렬

- 왼쪽 정렬
- 가운데 정렬
- 오른쪽 정렬
- 양쪽 정렬

### 📄 문서 관리

- 새 문서
- HTML 내보내기
- JSON 내보내기
- 전체화면 모드

## 4. 컴포넌트 사용법

### 기본 사용법

```vue
<template>
  <AdvancedTiptapEditor
    v-model="content"
    placeholder="내용을 입력하세요..."
    height="500px"
    @change="handleChange"
  />
</template>

<script setup>
import { ref } from 'vue';
import AdvancedTiptapEditor from 'src/components/AdvancedTiptapEditor.vue';

const content = ref('');

const handleChange = (value) => {
  console.log('에디터 내용:', value);
};
</script>
```

### Props

- `modelValue`: 에디터 내용 (v-model)
- `placeholder`: 플레이스홀더 텍스트
- `readonly`: 읽기 전용 모드
- `height`: 에디터 높이

### Events

- `update:modelValue`: 내용 변경 시
- `change`: 내용 변경 시 (HTML 문자열)

## 5. 메뉴바 기능

### 파일 메뉴

- **새 문서**: 내용 초기화
- **HTML 내보내기**: HTML 파일로 다운로드
- **JSON 내보내기**: JSON 파일로 다운로드

### 편집 메뉴

- **실행 취소**: Ctrl+Z
- **다시 실행**: Ctrl+Y
- **모두 선택**: Ctrl+A

### 보기 메뉴

- **전체화면**: 에디터를 전체화면으로 표시

## 6. 툴바 기능

### 텍스트 스타일

- 제목 1, 2, 3
- 굵게, 기울임, 밑줄, 취소선, 코드

### 색상 도구

- 텍스트 색상 팔레트
- 배경색 하이라이트 팔레트

### 정렬 도구

- 왼쪽, 가운데, 오른쪽, 양쪽 정렬

### 목록 도구

- 글머리 기호 목록
- 번호 매기기 목록
- 체크리스트

### 블록 도구

- 인용문
- 코드 블록

### 미디어 도구

- 링크 삽입
- 이미지 삽입
- 테이블 삽입

### 기타 도구

- 서식 지우기
- 실행 취소/다시 실행

## 7. 예제 페이지

`/sandbox/advanced-tiptap` 경로에서 고급 Tiptap 에디터 예제를 확인할 수 있습니다.

## 8. 스타일링

고급 Tiptap 에디터는 CSS로 커스터마이징할 수 있습니다:

```css
/* 에디터 컨테이너 */
.advanced-tiptap-editor {
  border: 1px solid #ddd;
  border-radius: 8px;
}

/* 메뉴바 */
.menubar {
  background-color: #f5f5f5;
  border-bottom: 1px solid #ddd;
}

/* 툴바 */
.toolbar {
  background-color: #f8f9fa;
  border-bottom: 1px solid #ddd;
}

/* 에디터 영역 */
:deep(.ProseMirror) {
  padding: 16px;
  outline: none;
}

/* 테이블 스타일 */
:deep(.ProseMirror table) {
  border-collapse: collapse;
  width: 100%;
}

:deep(.ProseMirror table th) {
  background-color: #f1f3f4;
  font-weight: bold;
}

/* 코드 블록 */
:deep(.ProseMirror pre) {
  background: #0d0d0d;
  color: #fff;
  border-radius: 0.5rem;
  padding: 0.75rem 1rem;
}

/* 체크리스트 */
:deep(.ProseMirror .task-list-item) {
  display: flex;
  align-items: flex-start;
}
```

## 9. 추가 확장 기능

필요에 따라 다음 확장 기능들을 추가할 수 있습니다:

```bash
# 협업 기능
npm install @tiptap/extension-collaboration @tiptap/extension-collaboration-cursor

# 멘션 기능
npm install @tiptap/extension-mention @tiptap/extension-mention-suggestion

# 수학 공식
npm install @tiptap/extension-math

# 다이어그램
npm install @tiptap/extension-diagram

# 드래그 앤 드롭
npm install @tiptap/extension-drag-handle

# 검색 및 바꾸기
npm install @tiptap/extension-find-and-replace
```

## 10. 문제 해결

### 모듈을 찾을 수 없는 오류

```bash
npm install @tiptap/vue-3 @tiptap/pm @tiptap/starter-kit
```

### 코드 블록 오류

```bash
npm install @tiptap/extension-code-block
```

### TypeScript 오류

```typescript
// tsconfig.json에 추가
{
  "compilerOptions": {
    "types": ["@tiptap/vue-3"]
  }
}
```

## 11. 성능 최적화

### 번들 크기 최적화

```javascript
// 필요한 확장 기능만 import
import StarterKit from '@tiptap/starter-kit';
import Table from '@tiptap/extension-table';
import TaskList from '@tiptap/extension-task-list';
```

### 메모리 최적화

```javascript
// 에디터 정리
onBeforeUnmount(() => {
  editor.value?.destroy();
});
```

## 12. 라이선스

Tiptap은 MIT 라이선스로 배포되므로 상업적 사용이 자유롭습니다.

| 에디터  | 라이선스 | 상업적 사용 | 무료 |
| ------- | -------- | ----------- | ---- |
| Tiptap  | MIT      | ✅          | ✅   |
| TinyMCE | GPL v2   | ⚠️ 제한적   | ✅   |
| Quill   | BSD      | ✅          | ✅   |
